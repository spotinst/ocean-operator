package controller

import (
	"github.com/go-logr/logr"
	"github.com/spotinst/ocean-operator/internal/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// NewRequestId returns random generated UUID.
func NewRequestId() string {
	return uuid.NewV4().String()
}

// NewRequestLogger returns a new logr.Logger with reconcile.Request's metadata.
func NewRequestLogger(logger logr.Logger, req reconcile.Request, reqID string) logr.Logger {
	return logger.WithValues(
		"Request.Id", reqID,
		"Request.Namespace", req.Namespace,
		"Request.Name", req.Name)
}
