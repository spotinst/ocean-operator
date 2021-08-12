// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controllers

import (
	"context"
	"fmt"
	"sort"

	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/pkg/tide"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newCondition creates a new OceanComponent condition.
func newCondition(condType oceanv1alpha1.OceanComponentConditionType,
	status corev1.ConditionStatus, reason, message string,
) *oceanv1alpha1.OceanComponentCondition {
	return &oceanv1alpha1.OceanComponentCondition{
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
	}
}

// newConditionf returns a new OceanComponent condition with arguments.
func newConditionf(condType oceanv1alpha1.OceanComponentConditionType,
	status corev1.ConditionStatus, reason, message string, args ...interface{},
) *oceanv1alpha1.OceanComponentCondition {
	return newCondition(condType, status, reason, fmt.Sprintf(message, args...))
}

// hasCondition returns true if the given status has the given condition.
func hasCondition(status *oceanv1alpha1.OceanComponentStatus,
	condition oceanv1alpha1.OceanComponentCondition) bool {
	c := status.Conditions
	for _, e := range c {
		if e.Type == condition.Type {
			return true
		}
	}
	return false
}

// setCondition updates the OceanComponent to include the provided
// condition. If the condition that we are about to add already exists and has
// the same status and reason then we are not going to update.
func setCondition(status *oceanv1alpha1.OceanComponentStatus,
	condition oceanv1alpha1.OceanComponentCondition) bool {
	currentCond := getCondition(*status, condition.Type)
	if currentCond != nil &&
		currentCond.Status == condition.Status &&
		currentCond.Reason == condition.Reason {
		return false
	}
	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := filterOutCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
	return true
}

// removeCondition removes the condition with the provided type.
func removeCondition(status *oceanv1alpha1.OceanComponentStatus,
	condType oceanv1alpha1.OceanComponentConditionType) {
	status.Conditions = filterOutCondition(status.Conditions, condType)
}

// getCurrentCondition returns the condition with the most recent
// update.
func getCurrentCondition(
	status oceanv1alpha1.OceanComponentStatus) *oceanv1alpha1.OceanComponentCondition {
	if len(status.Conditions) == 0 {
		return nil
	}
	sortMostRecent(&status)
	return &status.Conditions[0]
}

func sortMostRecent(status *oceanv1alpha1.OceanComponentStatus) {
	c := status.Conditions
	sort.Slice(c, func(i int, j int) bool {
		return c[i].LastUpdateTime.Time.After(c[j].LastUpdateTime.Time)
	})
	status.Conditions = c
}

// getCondition returns the condition with the provided type.
func getCondition(status oceanv1alpha1.OceanComponentStatus,
	condType oceanv1alpha1.OceanComponentConditionType) *oceanv1alpha1.OceanComponentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

// filterOutCondition returns a new slice of conditions without conditions with
// the provided type.
func filterOutCondition(conditions []oceanv1alpha1.OceanComponentCondition,
	condType oceanv1alpha1.OceanComponentConditionType) []oceanv1alpha1.OceanComponentCondition {
	var newConditions []oceanv1alpha1.OceanComponentCondition
	for _, condition := range conditions {
		if condition.Type == condType {
			continue
		}
		newConditions = append(newConditions, condition)
	}
	return newConditions
}

func getDeploymentConditions(ctx context.Context, client client.Client,
	objName types.NamespacedName) ([]*oceanv1alpha1.OceanComponentCondition, error) {
	conditions := make([]*oceanv1alpha1.OceanComponentCondition, 0)

	deployment := new(appsv1.Deployment)
	if err := client.Get(ctx, objName, deployment); err != nil {
		if apierrors.IsNotFound(err) {
			conditions = append(conditions, newCondition(
				oceanv1alpha1.OceanComponentConditionTypeAvailable,
				corev1.ConditionFalse,
				"DeploymentAbsent",
				"Deployment does not exist",
			))
			return conditions, nil // enough, return
		} else {
			return nil, err
		}
	}

	if deployment.Status.AvailableReplicas == 0 {
		conditions = append(conditions, newCondition(
			oceanv1alpha1.OceanComponentConditionTypeAvailable,
			corev1.ConditionFalse,
			"DeploymentUnavailable",
			"No pods are available",
		))
	} else {
		conditions = append(conditions, newCondition(
			oceanv1alpha1.OceanComponentConditionTypeAvailable,
			corev1.ConditionTrue,
			"DeploymentAvailable",
			"Pods are available",
		))
	}

	return conditions, nil
}

func getOceanControllerConditions(ctx context.Context, client client.Client,
	objName types.NamespacedName) ([]*oceanv1alpha1.OceanComponentCondition, error) {
	// ocean-controller's deployment name differs from its component name
	objName = types.NamespacedName{
		Namespace: metav1.NamespaceSystem,
		Name:      tide.LegacyOceanControllerDeployment,
	}
	return getDeploymentConditions(ctx, client, objName)
}

func getMetricsServerConditions(ctx context.Context, client client.Client,
	objName types.NamespacedName) ([]*oceanv1alpha1.OceanComponentCondition, error) {
	return getDeploymentConditions(ctx, client, objName)
}
