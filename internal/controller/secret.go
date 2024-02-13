package controller

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-password/password"
	unagexcomv1 "github.com/unagex/metabase-operator/api/v1"
	"github.com/unagex/metabase-operator/internal/controller/common"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (r *MetabaseReconciler) ManageSecret(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	sec := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: metabase.Namespace,
		Name:      metabase.Name,
	}, sec)

	// create if secret does not exist
	if k8serrors.IsNotFound(err) {
		sec, err := r.GetSecret(metabase)
		if err != nil {
			return err
		}

		err = r.Create(ctx, sec)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating secret: %w", err)
		}
		if err == nil {
			r.Log.Info("secret created")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting secret: %w", err)
	}

	return nil
}

func (r *MetabaseReconciler) GetSecret(metabase *unagexcomv1.Metabase) (*corev1.Secret, error) {
	ls := common.GetLabels(metabase.Name, "secret")
	password, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return nil, err
	}

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name + "-secret",
			Namespace: metabase.Namespace,
			Labels:    ls,
		},
		Immutable: ptr.To(true),
		Data: map[string][]byte{
			"PASSWORD": []byte(password),
		},
	}
	err = controllerruntime.SetControllerReference(metabase, sec, r.Scheme)
	if err != nil {
		return nil, err
	}

	return sec, nil
}
