package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsBeingDeleted returns whether this object has been requested to be deleted.
func IsBeingDeleted(obj metav1.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// HasFinalizer returns true if the given object has the given finalizer.
func HasFinalizer(o metav1.Object, finalizer string) bool {
	finalizers := o.GetFinalizers()
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

// AddFinalizer adds the given finalizer to the given object.
func AddFinalizer(obj metav1.Object, finalizer string) bool {
	finalizers := obj.GetFinalizers()
	controllerutil.AddFinalizer(obj, finalizer)
	return len(finalizers) < len(obj.GetFinalizers())
}

// RemoveFinalizer removes the given finalizer from the given object.
func RemoveFinalizer(obj metav1.Object, finalizer string) bool {
	finalizers := obj.GetFinalizers()
	controllerutil.RemoveFinalizer(obj, finalizer)
	return len(finalizers) > len(obj.GetFinalizers())
}
