package controller

import (
	"fmt"

	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewCondition returns a new StatusCondition.
func NewCondition(
	typ oceanv1.ConditionType,
	status corev1.ConditionStatus,
	reason oceanv1.ConditionReason,
	message string) oceanv1.StatusCondition {

	return oceanv1.StatusCondition{
		Type:               typ,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
		LastUpdateTime:     metav1.Now(),
	}
}

// NewConditionf returns a new StatusCondition with arguments.
func NewConditionf(
	typ oceanv1.ConditionType,
	status corev1.ConditionStatus,
	reason oceanv1.ConditionReason,
	message string,
	args ...interface{}) oceanv1.StatusCondition {

	return NewCondition(typ, status, reason, fmt.Sprintf(message, args...))
}

// HasCondition returns true if the given status has the given condition.
func HasCondition(status *oceanv1.ReconcileStatus, condition oceanv1.StatusCondition) bool {
	c := status.Conditions
	for _, e := range c {
		if e.Type == condition.Type {
			return true
		}
	}
	return false
}

// AddCondition adds/updates the given condition to the given status.
func AddCondition(status *oceanv1.ReconcileStatus, condition oceanv1.StatusCondition) {
	var existingCond *oceanv1.StatusCondition
	for i, cond := range status.Conditions {
		if cond.Type == condition.Type {
			// can't take a pointer to an iteration variable
			existingCond = &status.Conditions[i]
			break
		}
	}

	if existingCond == nil {
		status.Conditions = append(status.Conditions, condition)
		existingCond = &status.Conditions[len(status.Conditions)-1]
	}

	if existingCond.Status != condition.Status {
		existingCond.LastTransitionTime = metav1.Now()
	}

	existingCond.Status = condition.Status
	existingCond.Reason = condition.Reason
	existingCond.Message = condition.Message
}

// RemoveCondition removes the given condition from the given status.
func RemoveCondition(status *oceanv1.ReconcileStatus, condition oceanv1.StatusCondition) {
	c := status.Conditions
	for i, e := range c {
		if e.Type == condition.Type {
			c = append(c[:i], c[i+1:]...)
		}
	}
	status.Conditions = c
}
