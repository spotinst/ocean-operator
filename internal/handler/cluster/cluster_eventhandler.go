package cluster

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
// events for oceanv1.Cluster objects.
type EventHandler struct {
	// Logger allows the event handler to log messages.
	Log logr.Logger

	// Ocean represents an instance of the Ocean interface.
	Ocean spot.Ocean
}

// OnCreate is called in response to an create event.
func (x *EventHandler) OnCreate(ctx context.Context, evt event.CreateEvent) (runtime.Object, error) {
	x.Log.Info("Handling cluster creation")

	// Extract the underlying Cluster object.
	clusterObj, ok := evt.Object.(*oceanv1.Cluster)
	if !ok {
		return nil, handler.ErrNotImplemented
	}

	// Convert from object.
	cluster, err := x.Ocean.NewClusterConverter().FromObject(clusterObj)
	if err != nil {
		return nil, err
	}

	// Create the cluster.
	oceanCluster, err := x.Ocean.CreateCluster(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream cluster: %w", err)
	}

	// Convert to object.
	return x.Ocean.NewClusterConverter().ToObject(oceanCluster)
}

// OnUpdate is called in response to an update event.
func (x *EventHandler) OnUpdate(ctx context.Context, evt event.UpdateEvent) (runtime.Object, error) {
	x.Log.Info("Handling cluster configuration change")

	// Extract the underlying Cluster object.
	clusterObj, ok := evt.ObjectNew.(*oceanv1.Cluster)
	if !ok {
		return nil, handler.ErrNotImplemented
	}

	// Convert from object.
	cluster, err := x.Ocean.NewClusterConverter().FromObject(clusterObj)
	if err != nil {
		return nil, err
	}

	// Update the cluster.
	oceanCluster, err := x.Ocean.UpdateCluster(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to update upstream cluster: %w", err)
	}

	// Convert to object.
	return x.Ocean.NewClusterConverter().ToObject(oceanCluster)
}

// OnDelete is called in response to a delete event.
func (x *EventHandler) OnDelete(ctx context.Context, evt event.DeleteEvent) error {
	x.Log.Info("Handling cluster deletion")

	// If there is no upstream cluster, then let's bail early.
	if evt.Object == nil {
		return nil
	}

	// Extract the underlying Cluster object.
	clusterObj, ok := evt.Object.(*oceanv1.Cluster)
	if !ok {
		return handler.ErrNotImplemented
	}

	// Delete the upstream cluster.
	if oceanID := clusterObj.Status.OceanID; oceanID != "" {
		if err := x.Ocean.DeleteCluster(ctx, oceanID); err != nil {
			return fmt.Errorf("failed to delete upstream cluster: %w", err)
		}
	}

	return nil
}
