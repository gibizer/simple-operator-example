/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	testv1beta1 "github.com/gibizer/test-operator/api/v1beta1"
)

// ServiceWithDBReconciler reconciles a ServiceWithDB object
type ServiceWithDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.test.org,resources=servicewithdbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.test.org,resources=servicewithdbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.test.org,resources=servicewithdbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceWithDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceWithDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//return NewServiceWithDBReconcileReq(ctx, req, r.Client).Handle()
	return resultOK.Unwrap()
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceWithDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1beta1.ServiceWithDB{}).
		Complete(r)
}
