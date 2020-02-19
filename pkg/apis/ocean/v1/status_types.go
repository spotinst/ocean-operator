package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ReconcileStatus is the most recently observed status of a resource. It's used
// to communicate success or failure and the error message.
type ReconcileStatus struct {
	// The generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Represents the latest available observations of the current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []StatusCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// Extended data associated with the status. This field is optional.
	Details StatusDetails `json:"details,omitempty"`
}

// StatusCondition contains condition information.
type StatusCondition struct {
	// Type of the condition.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The reason for the condition's last transition.
	Reason ConditionReason `json:"reason,omitempty"`
}

// ConditionReason is an enumeration of possible failure causes.
type ConditionReason string

const (
	// The server has declined to indicate a specific reason.
	ConditionReasonUnknown ConditionReason = ""
)

// ConditionType is an enumeration of possible condition types.
type ConditionType string

const (
	// ConditionTypeReady indicates that the Cluster/LaunchSpec is ready and
	// able to scale if necessary: it's correctly configured and isn't disabled.
	ConditionTypeReady ConditionType = "Ready"
)

// StatusDetails is a set of additional properties that MAY be set by the
// controller to provide additional information.
type StatusDetails struct {
	// Unique ID of the reconcile request.
	ReconcileRequestUID types.UID `json:"reconcileRequestUid,omitempty"`
}
