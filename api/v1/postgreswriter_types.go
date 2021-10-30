/*
Copyright 2021.

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

//+kubebuilder:validation:Required
type PostgresWriterSpec struct {

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=string
	Table string `json:"table"`

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=integer
	//+kubebuilder:validation:Minimum=0
	Age int `json:"age"`

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=string
	Name string `json:"name"`

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=string
	Country string `json:"country"`
}

type WriteResult string

const (
	Success WriteResult = "success"
	Failed  WriteResult = "failed"
)

// PostgresWriterStatus defines the observed state of PostgresWriter
// +kubebuilder:printcolumn:JSONPath=".status.result", name=WriteResult, type=string
type PostgresWriterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:Type=string
	Result WriteResult `json:"result"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PostgresWriter is the Schema for the postgreswriters API
type PostgresWriter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresWriterSpec   `json:"spec,omitempty"`
	Status PostgresWriterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PostgresWriterList contains a list of PostgresWriter
type PostgresWriterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresWriter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresWriter{}, &PostgresWriterList{})
}
