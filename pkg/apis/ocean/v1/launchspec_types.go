package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LaunchSpecSpec defines the desired state of LaunchSpec
type LaunchSpecSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	OceanID            string              `json:"oceanId,omitempty"`
	ImageID            string              `json:"imageId,omitempty"`
	UserData           string              `json:"userData,omitempty"`
	RootVolumeSize     int                 `json:"rootVolumeSize,omitempty"`
	SecurityGroupIDs   []string            `json:"securityGroupIds,omitempty"`
	SubnetIDs          []string            `json:"subnetIds,omitempty"`
	IAMInstanceProfile *IAMInstanceProfile `json:"iamInstanceProfile,omitempty"`
	Labels             []*Label            `json:"labels,omitempty"`
	Taints             []*Taint            `json:"taints,omitempty"`
}

// LaunchSpecStatus defines the observed status of LaunchSpec.
type LaunchSpecStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Most recently observed status of the LaunchSpec.
	ReconcileStatus `json:",inline"`
	// Unique ID of the LaunchSpec.
	LaunchSpecID string `json:"launchSpecId,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LaunchSpec allows you to configure a workload type for your Cluster.
//
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=launchspecs,scope=Namespaced
// +kubebuilder:printcolumn:name="Spec ID",type="string",JSONPath=".status.launchSpecId",description="The unique ID of the LaunchSpec"
// +kubebuilder:printcolumn:name="Ocean ID",type="string",JSONPath=".spec.oceanId",description="The unique ID of the Cluster"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type LaunchSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the LaunchSpec.
	Spec LaunchSpecSpec `json:"spec,omitempty"`
	// Most recently observed status of the LaunchSpec.
	Status LaunchSpecStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LaunchSpecList contains a list of LaunchSpec.
type LaunchSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of LaunchSpecs.
	Items []LaunchSpec `json:"items"`
}

// Label defines the desired state of Label.
type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Taint defines the desired state of Taint.
type Taint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

func init() {
	SchemeBuilder.Register(&LaunchSpec{}, &LaunchSpecList{})
}
