package handler

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// ErrNotImplemented is the error returned if a method is not implemented.
var ErrNotImplemented = errors.New("handler: not implemented")

type (
	// EventHandler defines the interface that should be implemented by all
	// event handlers.
	EventHandler interface {
		// OnCreate is called in response to an create event.
		OnCreate(context.Context, event.CreateEvent) (runtime.Object, error)

		// OnUpdate is called in response to an update event.
		OnUpdate(context.Context, event.UpdateEvent) (runtime.Object, error)

		// OnDelete is called in response to a delete event.
		OnDelete(context.Context, event.DeleteEvent) error
	}

	// Event is an event describing a change that has happened to an object.
	Event interface{}
)

// EventHandlerFuncs implements EventHandler.
type EventHandlerFuncs struct {
	// CreateFunc is called in response to an create event.
	CreateFunc func(context.Context, event.CreateEvent) (runtime.Object, error)

	// UpdateFunc is called in response to an update event.
	UpdateFunc func(context.Context, event.UpdateEvent) (runtime.Object, error)

	// DeleteFunc is called in response to a delete event.
	DeleteFunc func(context.Context, event.DeleteEvent) error
}

// OnCreate implements EventHandler.
func (h EventHandlerFuncs) OnCreate(ctx context.Context, e event.CreateEvent) (runtime.Object, error) {
	if h.CreateFunc != nil {
		return h.CreateFunc(ctx, e)
	}

	return nil, ErrNotImplemented
}

// OnUpdate implements EventHandler.
func (h EventHandlerFuncs) OnUpdate(ctx context.Context, e event.UpdateEvent) (runtime.Object, error) {
	if h.UpdateFunc != nil {
		return h.UpdateFunc(ctx, e)
	}

	return nil, ErrNotImplemented
}

// OnDelete implements EventHandler.
func (h EventHandlerFuncs) OnDelete(ctx context.Context, e event.DeleteEvent) error {
	if h.DeleteFunc != nil {
		return h.DeleteFunc(ctx, e)
	}

	return ErrNotImplemented
}
