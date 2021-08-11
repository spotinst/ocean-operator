// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controller

import (
	"context"

	"github.com/spotinst/ocean-operator/pkg/log"
	ctrl "sigs.k8s.io/controller-runtime"
)

// RequestContext is an interface and serves as a base for contexts.
type RequestContext interface {
	context.Context

	// GetRequestId returns the request ID.
	GetRequestId() string
	// GetRequest returns the underlying request.
	GetRequest() ctrl.Request
	// GetLogger returns an initialized logger.
	GetLogger() log.Logger
	// Cancel cancels the request context and tells the reconciler to abandon its work.
	Cancel()
	// Canceled returns true whether the request context has been canceled.
	Canceled() bool
}

// reconcileContext implements the RequestContext interface.
type reconcileContext struct {
	context.Context

	requestID string
	request   ctrl.Request
	log       log.Logger

	// internal
	cancel context.CancelFunc // used for a cancellation signal
}

// NewRequestContext returns a new request context.
func NewRequestContext(ctx context.Context, req ctrl.Request,
	reqID string, reqLogger log.Logger) RequestContext {

	reqCtx, cancel := context.WithCancel(ctx)
	return &reconcileContext{
		Context:   reqCtx,
		request:   req,
		requestID: reqID,
		log:       reqLogger,
		cancel:    cancel,
	}
}

func (r *reconcileContext) GetRequestId() string     { return r.requestID }
func (r *reconcileContext) GetRequest() ctrl.Request { return r.request }
func (r *reconcileContext) GetLogger() log.Logger    { return r.log }

func (r *reconcileContext) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *reconcileContext) Canceled() bool {
	select {
	case <-r.Context.Done():
		return true
	default:
		return false
	}
}
