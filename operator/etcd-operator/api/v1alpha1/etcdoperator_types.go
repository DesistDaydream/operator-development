/*
Copyright 2020 DesistDaydream.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EtcdOperatorSpec defines the desired state of EtcdOperator
type EtcdOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size 是 etcd 集群节点的个数
	Size *int32 `json:"size"`
	// Image 是 etcd 所用镜像
	Image string `json:"image"`
}

// EtcdOperatorStatus defines the observed state of EtcdOperator
type EtcdOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// EtcdOperator is the Schema for the etcdoperators API
type EtcdOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdOperatorSpec   `json:"spec,omitempty"`
	Status EtcdOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdOperatorList contains a list of EtcdOperator
type EtcdOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdOperator{}, &EtcdOperatorList{})
}
