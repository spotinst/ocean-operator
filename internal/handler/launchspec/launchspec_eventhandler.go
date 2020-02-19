package launchspec

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spotinst/ocean-operator/internal/handler"
	"github.com/spotinst/ocean-operator/internal/spot"
	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// EventHandler implements the EventHandler interface by handling reconcile
// events for oceanv1.LaunchSpec objects.
type EventHandler struct {
	// Logger allows the event handler to log messages.
	Log logr.Logger

	// Ocean represents an instance of the Ocean interface.
	Ocean spot.Ocean
}

// OnCreate is called in response to an create event.
func (x *EventHandler) OnCreate(ctx context.Context, evt event.CreateEvent) (runtime.Object, error) {
	x.Log.Info("Handling launch spec creation")

	// Extract the underlying LaunchSpec object.
	specObj, ok := evt.Object.(*oceanv1.LaunchSpec)
	if !ok {
		return nil, handler.ErrNotImplemented
	}

	// Convert from object.
	spec, err := x.Ocean.NewLaunchSpecConverter().FromObject(specObj)
	if err != nil {
		return nil, err
	}

	// Create the spec.
	oceanLaunchSpec, err := x.Ocean.CreateLaunchSpec(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream launch spec: %w", err)
	}

	// Convert to object.
	return x.Ocean.NewLaunchSpecConverter().ToObject(oceanLaunchSpec)
}

// OnUpdate is called in response to an update event.
func (x *EventHandler) OnUpdate(ctx context.Context, evt event.UpdateEvent) (runtime.Object, error) {
	x.Log.Info("Handling launch spec configuration change")

	// Extract the underlying LaunchSpec object.
	specObj, ok := evt.ObjectNew.(*oceanv1.LaunchSpec)
	if !ok {
		return nil, handler.ErrNotImplemented
	}

	// Convert from object.
	spec, err := x.Ocean.NewLaunchSpecConverter().FromObject(specObj)
	if err != nil {
		return nil, err
	}

	// Update the spec.
	oceanLaunchSpec, err := x.Ocean.UpdateLaunchSpec(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to update upstream launch spec: %w", err)
	}

	// Convert to object.
	return x.Ocean.NewLaunchSpecConverter().ToObject(oceanLaunchSpec)
}

// OnDelete is called in response to a delete event.
func (x *EventHandler) OnDelete(ctx context.Context, evt event.DeleteEvent) error {
	x.Log.Info("Handling launch spec deletion")

	// If there is no upstream spec, then let's bail early.
	if evt.Object == nil {
		return nil
	}

	// Extract the underlying LaunchSpec object.
	specObj, ok := evt.Object.(*oceanv1.LaunchSpec)
	if !ok {
		return handler.ErrNotImplemented
	}

	// Delete the upstream spec.
	if specID := specObj.Status.LaunchSpecID; specID != "" {
		if err := x.Ocean.DeleteLaunchSpec(ctx, specID); err != nil {
			return fmt.Errorf("failed to delete upstream launch spec: %w", err)
		}
	}

	return nil
}
