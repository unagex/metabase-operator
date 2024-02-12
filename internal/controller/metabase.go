package controller

import (
	"context"
	"fmt"

	"github.com/unagex/metabase-operator/internal/controller/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	controllerruntime "sigs.k8s.io/controller-runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	unagexcomv1 "github.com/unagex/metabase-operator/api/v1"
)

func (r *MetabaseReconciler) ManageMetabase(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	dep := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: metabase.Namespace,
		Name:      metabase.Name,
	}, dep)

	// create if deployment does not exist
	if k8serrors.IsNotFound(err) {
		dep := r.GetDeployment(metabase)
		err := r.Create(ctx, dep)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating deployment: %w", err)
		}
		if err == nil {
			r.Log.Info("deployment created")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting deployment: %w", err)
	}

	return nil
}

func (r *MetabaseReconciler) GetDeployment(metabase *unagexcomv1.Metabase) *appsv1.Deployment {
	ls := common.GetLabels(metabase.Name)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name,
			Namespace: metabase.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "metabase/metabase:latest",
							Name:  "metabase",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 3000,
								},
							},
						},
					},
				},
			},
		},
	}
	_ = controllerruntime.SetControllerReference(metabase, dep, r.Scheme)

	return dep
}
