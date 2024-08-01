package kube

import (
	"context"
	"fmt"

	cranev1 "crane.cloud.cranom.tech/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateDeploymentFromApp(
	ctx context.Context,
	req ctrl.Request,
	application cranev1.Application,
	deploymentsClient *kubernetes.Clientset,
) error {
	applicationName := "application-" + req.Name
	deployment, err := deploymentsClient.AppsV1().Deployments(req.Namespace).Get(ctx, applicationName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// create deployment
			deploymentCFG := GetDeploymentCFGFromApp(application)
			_, err := deploymentsClient.AppsV1().Deployments(req.Namespace).Create(ctx, deploymentCFG, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("couldn't create deployment: %s", err)
			}
			return nil
		}
		return err
	}

	// update deployment
	deploymentCFG := GetDeploymentCFGFromApp(application)
	deployment.Spec = deploymentCFG.Spec
	_, err = deploymentsClient.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
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

/* func GetServiceCFGFromApp(
	app cranev1.Application,
) *appsv1.Deployment {

} */
