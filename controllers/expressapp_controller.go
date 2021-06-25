/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	resources "github.com/vincent-pli/serverless-operand/controllers/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	knativeserving "knative.dev/serving/pkg/apis/serving/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ibmdevv1alpha1 "github.com/vincent-pli/serverless-operand/api/v1alpha1"
)

// ExpressappReconciler reconciles a Expressapp object
type ExpressappReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ibm.dev.ibm.dev,resources=expressapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ibm.dev.ibm.dev,resources=expressapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ibm.dev.ibm.dev,resources=expressapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Expressapp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ExpressappReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("expressapp", req.NamespacedName)

	// Fetch the Memcached instance
	expressApp := &ibmdevv1alpha1.Expressapp{}
	err := r.Get(ctx, req.NamespacedName, expressApp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			r.Log.Info("expressApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Error(err, "Failed to get expressApp")
		return ctrl.Result{}, err
	}

	// Check if the ksvc already exists, if not create a new one
	found := &knativeserving.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-ksvc", expressApp.Name), Namespace: expressApp.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		ksvc := resources.MakeKSVC(expressApp)
		r.Log.Info("Creating a new KSVC", "KSVC.Namespace", ksvc.Namespace, "KSVC.Name", ksvc.Name)

		if err := controllerutil.SetControllerReference(expressApp, ksvc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		err = r.Create(ctx, ksvc)
		if err != nil {
			r.Log.Error(err, "Failed to create new KSVC", "KSVC.Namespace", ksvc.Namespace, "KSVC.Name", ksvc.Name)
			return ctrl.Result{}, err
		}
		// KSVC created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get KSVC")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExpressappReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ibmdevv1alpha1.Expressapp{}).
		Complete(r)
}
