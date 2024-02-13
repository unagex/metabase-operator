package controller

import (
	"context"
	"fmt"

	"github.com/unagex/metabase-operator/internal/controller/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	ls := common.GetLabels(metabase.Name, "metabase")
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
					Volumes: []corev1.Volume{
						{
							Name: metabase.Name + "-secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: metabase.Name + "-secret",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image:           metabase.Spec.Metabase.Image,
							ImagePullPolicy: metabase.Spec.Metabase.ImagePullPolicy,
							Name:            "metabase",
							Resources:       metabase.Spec.Metabase.Resources,
							StartupProbe: &corev1.Probe{
								// 80 * 10 = 800 seconds for the pod before being restarted
								FailureThreshold: 80,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/api/health",
										Port:   intstr.FromString("http"),
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold: 6,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/api/health",
										Port:   intstr.FromString("http"),
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/api/health",
										Port:   intstr.FromString("http"),
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 3000,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MB_DB_TYPE",
									Value: "postgres",
								},
								{
									Name:  "MB_DB_PORT",
									Value: "5432",
								},
								{
									Name:  "MB_DB_DBNAME",
									Value: "metabaseappdb",
								},
								{
									Name:  "MB_DB_USER",
									Value: "user",
								},
								{
									Name: "MB_DB_PASS",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: metabase.Name + "-secret",
											},
											Key: "PASSWORD",
										},
									},
								},
								{
									Name:  "MB_DB_HOST",
									Value: fmt.Sprintf("%s.%s.svc.cluster.local", metabase.Name+"-psql", metabase.Namespace),
								},
								{
									// TODO: make this var dynamic
									Name:  "JAVA_TIMEZONE",
									Value: "US/Pacific",
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
