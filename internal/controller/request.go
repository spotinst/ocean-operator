// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controller

import (
	uuid "github.com/satori/go.uuid"
	"github.com/spotinst/ocean-operator/pkg/log"
	ctrl "sigs.k8s.io/controller-runtime"
)

// NewRequestId returns random generated UUID.
func NewRequestId() string {
	return uuid.NewV4().String()
}

// NewRequestLog returns a new log.Logger with reconcile.Request's metadata.
func NewRequestLog(log log.Logger, req ctrl.Request, reqID string) log.Logger {
	return log.WithValues(
		"request", reqID,
		"namespace", req.Namespace,
		"name", req.Name)
}
