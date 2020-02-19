package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// reconcileBase implements the ReconcileBase interface.
type reconcileBase struct {
	// This Client, initialized using mgr.Client(), is a split Client that reads
	// objects from the cache and writes to the apiserver.
	Client   client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// NewReconcile returns a new base Reconciler.
func NewReconcile(client client.Client, scheme *runtime.Scheme,
	recorder record.EventRecorder) ReconcileBase {

	return &reconcileBase{
		Client:   client,
		Scheme:   scheme,
		Recorder: recorder,
	}
}

func (r *reconcileBase) GetClient() client.Client          { return r.Client }
func (r *reconcileBase) GetRecorder() record.EventRecorder { return r.Recorder }
func (r *reconcileBase) GetScheme() *runtime.Scheme        { return r.Scheme }

func (r *reconcileBase) CreateOrUpdateResource(ctx context.Context, obj runtime.Object) error {
	noopMutateFn := controllerutil.MutateFn(func() error { return nil })

	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, obj, noopMutateFn); err != nil {
		err = fmt.Errorf("failed to create/update object: %w", err)
		r.RecordEvent(obj, corev1.EventTypeWarning, "FailedCreateOrUpdate", err.Error())
		return err
	}

	return nil
}

func (r *reconcileBase) CreateOrUpdateResources(ctx context.Context, objs []runtime.Object) error {
	for _, obj := range objs {
		if err := r.CreateOrUpdateResource(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (r *reconcileBase) UpdateResource(ctx context.Context, obj runtime.Object) error {
	if err := r.Client.Update(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
		err = fmt.Errorf("failed to update object: %w", err)
		r.RecordEvent(obj, corev1.EventTypeWarning, "FailedUpdate", err.Error())
		return err
	}

	return nil
}

func (r *reconcileBase) UpdateResources(ctx context.Context, objs []runtime.Object) error {
	for _, obj := range objs {
		if err := r.UpdateResource(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (r *reconcileBase) DeleteResource(ctx context.Context, obj runtime.Object) error {
	if err := r.Client.Delete(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
		err = fmt.Errorf("failed to delete object: %w", err)
		r.RecordEvent(obj, corev1.EventTypeWarning, "FailedDelete", err.Error())
		return err
	}

	return nil
}

func (r *reconcileBase) DeleteResources(ctx context.Context, objs []runtime.Object) error {
	for _, obj := range objs {
		if err := r.DeleteResource(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (r *reconcileBase) UpdateStatus(ctx context.Context, obj runtime.Object) error {
	if err := r.Client.Status().Update(ctx, obj); err != nil {
		err = fmt.Errorf("failed to update status: %w", err)
		r.RecordEvent(obj, corev1.EventTypeWarning, "FailedUpdateStatus", err.Error())
		return err
	}

	return nil
}

func (r *reconcileBase) RecordEvent(obj runtime.Object, typ, reason, message string) {
	r.Recorder.Event(obj, typ, reason, message)
}

func (r *reconcileBase) RecordEventf(obj runtime.Object, typ, reason, message string, args ...interface{}) {
	r.Recorder.Eventf(obj, typ, reason, message, args...)
}
