package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	controllerruntime "sigs.k8s.io/controller-runtime"

	unagexcomv1 "github.com/unagex/metabase-operator/api/v1"
	"github.com/unagex/metabase-operator/internal/controller/common"
)

func (r *MetabaseReconciler) ManageServices(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	err := r.ManageServiceHTTP(ctx, metabase)
	if err != nil {
		return err
	}

	err = r.ManageServiceDB(ctx, metabase)
	return err

}

func (r *MetabaseReconciler) ManageServiceHTTP(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: metabase.Namespace,
		Name:      metabase.Name + "-http",
	}, svc)

	// create if service does not exist
	if k8serrors.IsNotFound(err) {
		svc = r.GetServiceHTTP(metabase)
		err = r.Create(ctx, svc)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating service http: %w", err)
		}
		if err == nil {
			r.Log.Info("service http created")
		}
		return nil
	}

	return nil
}

func (r *MetabaseReconciler) ManageServiceDB(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: metabase.Namespace,
		Name:      metabase.Name + "-psql",
	}, svc)

	// create if service does not exist
	if k8serrors.IsNotFound(err) {
		svc = r.GetServiceDB(metabase)
		err = r.Create(ctx, svc)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating service db: %w", err)
		}
		if err == nil {
			r.Log.Info("service db created")
		}
		return nil
	}

	return nil
}

func (r *MetabaseReconciler) GetServiceHTTP(metabase *unagexcomv1.Metabase) *corev1.Service {
	ls := common.GetLabels(metabase.Name, "metabase")
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name + "-http",
			Namespace: metabase.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					TargetPort: intstr.FromString("http"),
					Port:       3000,
				},
			},
		},
	}
	_ = controllerruntime.SetControllerReference(metabase, svc, r.Scheme)

	return svc
}

func (r *MetabaseReconciler) GetServiceDB(metabase *unagexcomv1.Metabase) *corev1.Service {
	ls := common.GetLabels(metabase.Name, "database")
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name + "-psql",
			Namespace: metabase.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Name:       "psql",
					TargetPort: intstr.FromString("psql"),
					Port:       5432,
				},
			},
		},
	}
	_ = controllerruntime.SetControllerReference(metabase, svc, r.Scheme)

	return svc
}
