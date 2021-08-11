// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsBeingDeleted returns whether this object has been requested to be deleted.
func IsBeingDeleted(obj client.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// HasFinalizer returns true if the given object has the given finalizer.
func HasFinalizer(o client.Object, finalizer string) bool {
	finalizers := o.GetFinalizers()
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

// AddFinalizer adds the given finalizer to the given object.
func AddFinalizer(obj client.Object, finalizer string) bool {
	finalizers := obj.GetFinalizers()
	controllerutil.AddFinalizer(obj, finalizer)
	return len(finalizers) < len(obj.GetFinalizers())
}

// RemoveFinalizer removes the given finalizer from the given object.
func RemoveFinalizer(obj client.Object, finalizer string) bool {
	finalizers := obj.GetFinalizers()
	controllerutil.RemoveFinalizer(obj, finalizer)
	return len(finalizers) > len(obj.GetFinalizers())
}
