/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/go-logr/logr"
	"github.com/juju/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	secretsv1alpha1 "com.mithung.dev/supersecret/api/v1alpha1"
)

var KubeClient *kubernetes.Clientset

// SuperSecretReconciler reconciles a SuperSecret object
type SuperSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	// Metrics *MetricRecorder
}

//+kubebuilder:rbac:groups=secrets.com.mithung.dev,resources=supersecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secrets.com.mithung.dev,resources=supersecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secrets.com.mithung.dev,resources=supersecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SuperSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *SuperSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	fmt.Println("hi from controller")

	// TODO(user): your logic here

	// Handle panic and recover
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	// log := r.Log.WithValues("SuperSecret", req.NamespacedName)
	// // fmt.Println(log)
	// log.V(0).Info("==========start of the log==================")
	instance := &secretsv1alpha1.SuperSecret{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			return ctrl.Result{}, nil
		}

		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}

	// if no phase set, default to Pending
	if instance.Status.Phase == "" {
		println("setting phase of the object")
		instance.Status.Phase = "doingsomethingbuddy"
		println(instance.Status.Phase)
	}

	// update status
	err = r.Status().Update(context.TODO(), instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	}
	if config == nil {
		config = &rest.Config{}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.Sname,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"myvalue": []byte(instance.Spec.Ssecret),
		},
		Type: corev1.SecretTypeOpaque,
	}

	result, err := client.CoreV1().Secrets(instance.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Log.Error(err, "error gound")
	}
	fmt.Println("now created secret", result.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SuperSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.SuperSecret{}).
		Complete(r)
}
