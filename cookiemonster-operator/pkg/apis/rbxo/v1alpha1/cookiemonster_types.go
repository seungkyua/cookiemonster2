package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CookiemonsterSpec defines the desired state of Cookiemonster
// +k8s:openapi-gen=true
type Resource struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Target int64 `json:"target"`
}

type Namespace struct {
	Name     string `json:"name"`
	Resource []Resource `json:"resource"`
}

type Data struct{
	Namespace []Namespace `json:"namespace"`
	Size int32 `json:"size"`
	Interval int32 `json:"interval"`
	Duration int64 `json:"duration"`
	Slack 	bool `json:"slack"`
	Slackwebhook string `json:"slackwebhook"`
	Change bool `json:"change"`
	Bmcad string `json:"bmcad"`
}

type CookiemonsterSpec struct {
	Data Data `json:"data"`


	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// CookiemonsterStatus defines the observed state of Cookiemonster
// +k8s:openapi-gen=true
type CookiemonsterStatus struct {
	Nodes []string `json:"nodes"`
	Maps []string `json"maps"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cookiemonster is the Schema for the cookiemonsters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=cookiemonsters,scope=Namespaced
type Cookiemonster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CookiemonsterSpec   `json:"spec,omitempty"`
	Status CookiemonsterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CookiemonsterList contains a list of Cookiemonster
type CookiemonsterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cookiemonster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cookiemonster{}, &CookiemonsterList{})
}
