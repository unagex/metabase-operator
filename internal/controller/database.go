package controller

import (
	"context"
	"fmt"

	unagexcomv1 "github.com/unagex/metabase-operator/api/v1"
	"github.com/unagex/metabase-operator/internal/controller/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (r *MetabaseReconciler) ManageDatabase(ctx context.Context, metabase *unagexcomv1.Metabase) error {
	sts := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: metabase.Namespace,
		Name:      metabase.Name,
	}, sts)

	// create if statefulset does not exist
	if k8serrors.IsNotFound(err) {
		dep := r.GetStatefulSet(metabase)
		err := r.Create(ctx, dep)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating statefulset: %w", err)
		}
		if err == nil {
			r.Log.Info("statefulset created")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting statefulset: %w", err)
	}

	return nil
}

func (r *MetabaseReconciler) GetStatefulSet(metabase *unagexcomv1.Metabase) *appsv1.StatefulSet {
	ls := common.GetLabels(metabase.Name, "database")
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name,
			Namespace: metabase.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.StatefulSetSpec{
			// Only one replicas of metabase is allowed at the moment.
			Replicas: metabase.Spec.DB.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			ServiceName: metabase.Name + "-db",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: metabase.Name + "-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: metabase.Name,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image:           metabase.Spec.DB.Image,
							ImagePullPolicy: metabase.Spec.DB.ImagePullPolicy,
							Name:            metabase.Name,
							Ports: []corev1.ContainerPort{
								{
									Name:          "psql",
									ContainerPort: 5432,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh", "-c", fmt.Sprintf("pg_isready -U %s -d %s -q", "user", "metabaseappdb"),
										},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold: 6,
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh", "-c", fmt.Sprintf("pg_isready -U %s -d %s -q", "user", "metabaseappdb"),
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_USER",
									Value: "user",
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "password",
								},
								{
									Name:  "POSTGRES_DB",
									Value: "metabaseappdb",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/var/lib/postgresql/data",
									Name:      metabase.Name + "-storage",
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: r.getVCTs(metabase),
		},
	}
	_ = controllerruntime.SetControllerReference(metabase, sts, r.Scheme)

	return sts
}

func (r *MetabaseReconciler) getVCTs(metabase *unagexcomv1.Metabase) []corev1.PersistentVolumeClaim {
	r.Log.Info("the value is")
	r.Log.Info(fmt.Sprintf("%+v", metabase.Spec))
	vct := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metabase.Name + "-storage",
			Namespace: metabase.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: metabase.Spec.DB.Volume.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(metabase.Spec.DB.Volume.Size),
				},
			},
		},
	}
	_ = controllerruntime.SetControllerReference(metabase, vct, r.Scheme)

	return []corev1.PersistentVolumeClaim{*vct}
}
