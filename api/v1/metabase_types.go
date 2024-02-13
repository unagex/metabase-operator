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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MetabaseSpec defines the desired state of Metabase
type MetabaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default="postgres:latest"
	// +kubebuilder:validation:Optional
	Test string `json:"test"`

	// +kubebuilder:validation:Required
	DB DBSpec `json:"db"`
}

type DBSpec struct {
	// The image name to use for PostgreSQL containers.
	// +kubebuilder:default="postgres:latest"
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty"`

	// ImagePullPolicy is used to determine when Kubernetes will attempt to
	// pull (download) container images.
	// +kubebuilder:validation:Enum={Always,Never,IfNotPresent}
	// +kubebuilder:default="IfNotPresent"
	// +kubebuilder:validation:Optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Number of desired PostgreSQL pods.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	// +kubebuilder:validation:Optional
	Replicas *int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Required
	Volume VolumeSpec `json:"volume"`
}

type VolumeSpec struct {
	// StorageClassName defined for the volume.
	// +kubebuilder:validation:Optional
	StorageClassName *string `json:"storageClassName,omitempty"`

	// Size of the volume.
	// +kubebuilder:validation:default=10Gi
	// +kubebuilder:validation:Pattern=`^\d+(Gi|Gb|Ki|)$`
	// +kubebuilder:validation:Pattern=`^\d+(Ki|Mi|Gi|Ti|Pi|Ei|m|k|M|G|T|P|E)$`
	// +kubebuilder:validation:Required
	Size string `json:"size"`
}

// MetabaseStatus defines the observed state of Metabase
type MetabaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Metabase is the Schema for the metabases API
type Metabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   MetabaseSpec   `json:"spec"`
	Status MetabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MetabaseList contains a list of Metabase
type MetabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Metabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Metabase{}, &MetabaseList{})
}
