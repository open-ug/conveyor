package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cranev1 "crane.cloud.cranom.tech/api/v1"
	craneKubeUtils "crane.cloud.cranom.tech/kube"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(cranev1.AddToScheme(scheme))
}

type Reconciler struct {
	client.Client
	scheme     *runtime.Scheme
	kubeClient *kubernetes.Clientset
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("application", req.NamespacedName)
	log.Info("reconciling application")

	// create deployment if not exists
	deploymentsClient := r.kubeClient.AppsV1().Deployments(req.Namespace)

	var application cranev1.Application
	err := r.Client.Get(ctx, req.NamespacedName, &application)
	if err != nil {
		if k8serrors.IsNotFound(err) { // application not found, we can delete the resources
			applicationName := "application-" + req.Name
			err = deploymentsClient.Delete(ctx, applicationName, metav1.DeleteOptions{})
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("couldn't delete deployment: %s", err)
			}
		}
	}

	// create or update deployment
	err = craneKubeUtils.CreateDeploymentFromApp(ctx, req, application, r.kubeClient)

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("couldn't create or update deployment: %s", err)
	}

	return ctrl.Result{}, nil
}

func main() {
	var (
		config *rest.Config
		err    error
	)
	kubeconfigFilePath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(kubeconfigFilePath); errors.Is(err, os.ErrNotExist) { // if kube config doesn't exist, try incluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFilePath)
		if err != nil {
			panic(err.Error())
		}
	}

	// kubernetes client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctrl.SetLogger(zap.New())

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&cranev1.Application{}).
		Complete(&Reconciler{
			Client:     mgr.GetClient(),
			scheme:     mgr.GetScheme(),
			kubeClient: clientset,
		})

	if err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "error running manager")
		os.Exit(1)
	}

}

func int32Ptr(i int32) *int32 { return &i }
