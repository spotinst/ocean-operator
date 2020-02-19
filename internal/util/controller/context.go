package controller

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileContextBase implements the RequestContextBase interface.
type reconcileContextBase struct {
	context.Context

	RequestID string
	Request   reconcile.Request
	Log       logr.Logger

	// internal
	cancel context.CancelFunc // used for a cancellation signal
}

// NewRequestContext returns a new request context.
func NewRequestContext(ctx context.Context, req reconcile.Request,
	reqID string, reqLogger logr.Logger) RequestContextBase {

	reqCtx, cancel := context.WithCancel(ctx)
	return &reconcileContextBase{
		Context:   reqCtx,
		Request:   req,
		RequestID: reqID,
		Log:       reqLogger,
		cancel:    cancel,
	}
}

func (r *reconcileContextBase) GetRequestId() string          { return r.RequestID }
func (r *reconcileContextBase) GetRequest() reconcile.Request { return r.Request }
func (r *reconcileContextBase) GetLogger() logr.Logger        { return r.Log }

// Cancel cancels the request context and tells the ReconcileBase to abandon its work.
func (r *reconcileContextBase) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}

// Canceled returns true whether the request context has been canceled.
func (r *reconcileContextBase) Canceled() bool {
	select {
	case <-r.Context.Done():
		return true
	default:
		return false
	}
}
