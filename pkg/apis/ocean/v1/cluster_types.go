package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSpec defines the desired state of Cluster.
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file.
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Region     string      `json:"region,omitempty"`
	Compute    Compute     `json:"compute,omitempty"`
	Strategy   Strategy    `json:"strategy,omitempty"`
	Capacity   *Capacity   `json:"capacity,omitempty"`
	AutoScaler *AutoScaler `json:"autoScaler,omitempty"`
}

// ClusterStatus defines the observed status of Cluster.
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file.
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Most recently observed status of the Cluster.
	ReconcileStatus `json:",inline"`
	// Unique ID of the Cluster.
	OceanID string `json:"oceanId,omitempty"`
	// Observed status of the Cluster's nodes.
	Nodes NodesStatus `json:"nodes,omitempty"`
}

// NodesStatus defines the state of Nodes.
type NodesStatus struct {
	// Current number of nodes created by the Cluster.
	Current int32 `json:"current"`
	// Lower limit for the number of nodes that can be created by the Cluster.
	Minimum int32 `json:"minimum"`
	// Upper limit for the number of nodes that can be created by the Cluster.
	Maximum int32 `json:"maximum"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster consists of one or more Launch Specs and automatically scales your infrastructure.
//
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters,scope=Namespaced
// +kubebuilder:printcolumn:name="Cluster ID",type="string",JSONPath=".status.oceanId",description="The unique ID of the Cluster"
// +kubebuilder:printcolumn:name="Region",type="string",JSONPath=".spec.region",description="The region of the Cluster"
// +kubebuilder:printcolumn:name="Nodes",type="integer",JSONPath=".status.nodes.current",description="The current number of nodes created by the Cluster"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Cluster.
	Spec ClusterSpec `json:"spec,omitempty"`
	// Most recently observed status of the Cluster.
	Status ClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Clusters.
	Items []Cluster `json:"items"`
}

// Strategy defines the desired state of Strategy.
type Strategy struct {
	SpotPercentage           int32 `json:"spotPercentage,omitempty"`
	UtilizeReservedInstances bool  `json:"utilizeReservedInstances,omitempty"`
	FallbackToOnDemand       bool  `json:"fallbackToOd,omitempty"`
	DrainingTimeout          int32 `json:"drainingTimeout,omitempty"`
}

// Capacity defines the desired state of Capacity.
type Capacity struct {
	Minimum int32 `json:"minimum,omitempty"`
	Maximum int32 `json:"maximum,omitempty"`
	Target  int32 `json:"target,omitempty"`
}

// Compute defines the desired state of Compute.
type Compute struct {
	LaunchSpecification LaunchSpecification `json:"launchSpecification,omitempty"`
	SubnetIDs           []string            `json:"subnetIds,omitempty"`
}

// LaunchSpecification defines the desired state of LaunchSpecification.
type LaunchSpecification struct {
	AssociatePublicIPAddress bool                `json:"associatePublicIpAddress,omitempty"`
	ImageID                  string              `json:"imageId,omitempty"`
	KeyPair                  string              `json:"keyPair,omitempty"`
	UserData                 string              `json:"userData,omitempty"`
	RootVolumeSize           int32               `json:"rootVolumeSize,omitempty"`
	Monitoring               bool                `json:"monitoring,omitempty"`
	EBSOptimized             bool                `json:"ebsOptimized,omitempty"`
	SecurityGroupIDs         []string            `json:"securityGroupIds,omitempty"`
	IAMInstanceProfile       *IAMInstanceProfile `json:"iamInstanceProfile,omitempty"`
	Tags                     []*Tag              `json:"tags,omitempty"`
}

// IAMInstanceProfile defines the desired state of IAMInstanceProfile.
type IAMInstanceProfile struct {
	ARN  string `json:"arn,omitempty"`
	Name string `json:"name,omitempty"`
}

// Tag defines the desired state of Tag.
type Tag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// AutoScaler defines the desired state of AutoScaler.
type AutoScaler struct {
	IsEnabled      bool                      `json:"isEnabled,omitempty"`
	IsAutoConfig   bool                      `json:"isAutoConfig,omitempty"`
	Cooldown       int                       `json:"cooldown,omitempty"`
	Headroom       *AutoScalerHeadroom       `json:"headroom,omitempty"`
	ResourceLimits *AutoScalerResourceLimits `json:"resourceLimits,omitempty"`
	Down           *AutoScalerDown           `json:"down,omitempty"`
}

// AutoScalerHeadroom defines the desired state of AutoScalerHeadroom.
type AutoScalerHeadroom struct {
	CPUPerUnit    int `json:"cpuPerUnit,omitempty"`
	GPUPerUnit    int `json:"gpuPerUnit,omitempty"`
	MemoryPerUnit int `json:"memoryPerUnit,omitempty"`
	NumOfUnits    int `json:"numOfUnits,omitempty"`
}

// AutoScalerResourceLimits defines the desired state of AutoScalerResourceLimits.
type AutoScalerResourceLimits struct {
	MaxVCPU      int `json:"maxVCpu,omitempty"`
	MaxMemoryGiB int `json:"maxMemoryGib,omitempty"`
}

// AutoScalerDown defines the desired state of AutoScalerDown.
type AutoScalerDown struct {
	EvaluationPeriods      int `json:"evaluationPeriods,omitempty"`
	MaxScaleDownPercentage int `json:"maxScaleDownPercentage,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
