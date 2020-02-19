package controller

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ErrNotImplemented is the error returned if a method is not implemented.
var ErrNotImplemented = errors.New("controller: not implemented")

type (
	// RequestContextBase is an interface and serves as a base for Contexts.
	RequestContextBase interface {
		context.Context

		// GetRequestId returns the request ID.
		GetRequestId() string

		// GetRequest returns the underlying Request.
		GetRequest() reconcile.Request

		// GetLogger returns an initialized Logger.
		GetLogger() logr.Logger

		// Cancel cancels the request context and tells the Reconciler to
		// abandon its work.
		Cancel()

		// Canceled returns true whether the request context has been canceled.
		Canceled() bool
	}

	// ReconcileBase is an interface and serves as a base for Reconcilers.
	ReconcileBase interface {
		// GetClient returns an initialized Client. Defaults to DelegatingClient
		// that will use the cache for reads and the client for writes.
		GetClient() client.Client

		// GetScheme returns an initialized Scheme.
		GetScheme() *runtime.Scheme

		// GetEventRecorderFor returns an initialized EventRecorder.
		GetRecorder() record.EventRecorder

		// CreateOrUpdateResource creates a resource if it doesn't exist, and
		// updates (overwrites it), if it exists.
		CreateOrUpdateResource(ctx context.Context, obj runtime.Object) error

		// CreateOrUpdateResources creates one or more resources if they don't
		// exist, and updates (overwrites them), if they exist.
		CreateOrUpdateResources(ctx context.Context, objs []runtime.Object) error

		// UpdateResource updates an existing resource.
		// It doesn't fail if the resource doesn't exist.
		UpdateResource(ctx context.Context, obj runtime.Object) error

		// UpdateResources updates one ore more existing resources.
		// It doesn't fail if the resources don't exist.
		UpdateResources(ctx context.Context, objs []runtime.Object) error

		// DeleteResource deletes an existing resource.
		// It doesn't fail if the resource doesn't exist.
		DeleteResource(ctx context.Context, obj runtime.Object) error

		// DeleteResources deletes one ore more existing resources.
		// It doesn't fail if the resources don't exist.
		DeleteResources(ctx context.Context, objs []runtime.Object) error

		// UpdateStatus updates the fields corresponding to the status
		// sub-resource for the given obj.
		UpdateStatus(ctx context.Context, obj runtime.Object) error

		// RecordEvent constructs an event from the given information and puts
		// it in the queue for sending.
		RecordEvent(obj runtime.Object, typ, reason, message string)

		// RecordEventf is just like RecordEvent, but with Sprintf for the
		// message field.
		RecordEventf(obj runtime.Object, typ, reason, message string, args ...interface{})
	}
)
