/*
Copyright 2023.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ContainerSpec struct {
	Image string `json:"image"`
	Port  int32  `json:"port"`
}
type ServiceSpec struct {
	ServiceName     string `json:"serviceName,omitempty"`
	ServiceType     string `json:"serviceType"`
	ServicePort     int32  `json:"servicePort"`
	ServiceNodePort int32  `json:"serviceNodePort,omitempty"`
}
type CustomRSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentName string        `json:"deploymentName,omitempty"`
	Replicas       *int32        `json:"replicas"`
	Container      ContainerSpec `json:"container"`
	Service        ServiceSpec   `json:"service"`
}
type CustomRStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvailableReplicas int32 `json:"available_replicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CustomR is the Schema for the customrs API
type CustomR struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomRSpec   `json:"spec,omitempty"`
	Status CustomRStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CustomRList contains a list of CustomR
type CustomRList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomR `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomR{}, &CustomRList{})
}
