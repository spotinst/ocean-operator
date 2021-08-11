// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controller

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func Requeue(requeue bool) (ctrl.Result, error) {
	return ctrl.Result{Requeue: requeue}, nil
}

func RequeueError(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func RequeueAfterError(interval time.Duration, err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: interval}, err
}

func RequeueAfter(interval time.Duration) (ctrl.Result, error) {
	return RequeueAfterError(interval, nil)
}

func RequeueImmediately() (ctrl.Result, error) {
	return Requeue(true)
}

func NoRequeue() (ctrl.Result, error) {
	return RequeueError(nil)
}
