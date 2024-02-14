/*
Copyright 2024.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	unagexcomv1 "github.com/unagex/metabase-operator/api/v1"
)

// MetabaseReconciler reconciles a Metabase object
type MetabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

func (r *MetabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log = log.FromContext(ctx).WithName("Reconciler")

	metabase := &unagexcomv1.Metabase{}
	err := r.Get(ctx, req.NamespacedName, metabase)
	if k8serrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting metabase cr: %w", err)
	}

	err = r.ManageSecret(ctx, metabase)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ManageDatabase(ctx, metabase)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ManageMetabase(ctx, metabase)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ManageServices(ctx, metabase)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// filter to requeue when a dependant resource is created/updated/deleted.
	// secrets are the exception because we cannot recreate it with the same password.
	filter := handler.EnqueueRequestsFromMapFunc(func(_ context.Context, o client.Object) []reconcile.Request {
		ls := o.GetLabels()
		if ls["app.kubernetes.io/managed-by"] != "metabase-operator" {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Namespace: o.GetNamespace(),
					Name:      ls["app.kubernetes.io/instance"],
				},
			},
		}
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&unagexcomv1.Metabase{}).
		Watches(&appsv1.StatefulSet{}, filter).
		Watches(&appsv1.Deployment{}, filter).
		Watches(&corev1.Service{}, filter).
		Complete(r)
}
