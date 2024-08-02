package kube

import (
	"context"
	"fmt"
	"log"

	cranev1 "crane.cloud.cranom.tech/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func ApplyApplication(
	ctx context.Context,
	req ctrl.Request,
	application cranev1.Application,
	kubeClient *kubernetes.Clientset,
) error {
	created := false
	applicationName := "application-" + req.Name
	// check if deployment exists
	deployment, err := kubeClient.AppsV1().Deployments(req.Namespace).Get(ctx, applicationName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// create deployment
			deploymentCFG := GetDeploymentCFGFromApp(application)
			_, err := kubeClient.AppsV1().Deployments(req.Namespace).Create(ctx, deploymentCFG, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("couldn't create deployment: %s", err)
			}

		} else {
			return err
		}

	}

	// Check if service exists
	service, err := kubeClient.CoreV1().Services(req.Namespace).Get(ctx, applicationName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// create service
			serviceCFG := GetServiceCFGFromApp(application)
			_, err := kubeClient.CoreV1().Services(req.Namespace).Create(ctx, serviceCFG, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("couldn't create service: %s", err)
			}
			created = true
		} else {
			return err
		}
	}

	if created {
		return nil
	}

	// update service
	serviceCFG := GetServiceCFGFromApp(application)
	service.Spec = serviceCFG.Spec
	_, err = kubeClient.CoreV1().Services(req.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("couldn't update service: %s", err)
	}

	// update deployment
	deploymentCFG := GetDeploymentCFGFromApp(application)
	deployment.Spec = deploymentCFG.Spec
	_, err = kubeClient.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("couldn't update deployment: %s", err)
	}

	return nil
}

func GetDeploymentCFGFromApp(
	app cranev1.Application,
) *appsv1.Deployment {

	// Loop through app.Spec.Ports and add them to the container
	ports := []corev1.ContainerPort{}
	for _, port := range app.Spec.Ports {
		ports = append(ports, corev1.ContainerPort{
			Name:          port.Domain,
			ContainerPort: int32(port.Internal),
		})
	}

	// generate envFrom
	envFrom := []corev1.EnvFromSource{}
	envFrom = append(envFrom, corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: app.Spec.EnvFrom,
			},
		},
	})

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "application-" + app.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "application-" + app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "application-" + app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "application",
							Image:   app.Spec.Image,
							Ports:   ports,
							EnvFrom: envFrom,
						},
					},
				},
			},
		},
	}
	return deployment
}

func GetServiceCFGFromApp(
	app cranev1.Application,
) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "application-" + app.Name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "application-" + app.Name,
			},
			Ports: []corev1.ServicePort{},
			Type:  corev1.ServiceTypeNodePort,
		},
	}

	for _, port := range app.Spec.Ports {
		service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
			Name:     port.Domain,
			Port:     int32(port.External),
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(port.Internal),
			},
		})
	}

	return service
}

func GetPersistentVolumeClaimFromApp(
	app cranev1.Application,
) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "application-" + app.Name,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceQuantity(app.Spec.Resources.Storage),
				},
			},
		},
	}

	return pvc
}

func resourceQuantity(i int) resource.Quantity {
	return *resource.NewQuantity(int64(i), resource.BinarySI)
}

func DeleteApplication(
	ctx context.Context,
	req ctrl.Request,
	kubeClient *kubernetes.Clientset,
) error {
	applicationName := "application-" + req.Name
	// check if deployment exists
	err := kubeClient.AppsV1().Deployments(req.Namespace).Delete(ctx, applicationName, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Check if service exists
	err = kubeClient.CoreV1().Services(req.Namespace).Delete(ctx, applicationName, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Printf("service not found")
			return nil
		} else {
			return err
		}
	}

	return nil
}
