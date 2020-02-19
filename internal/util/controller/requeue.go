package controller

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func RequeueIfError(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func RequeueAfter(interval time.Duration, err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: interval}, err
}

func RequeueImmediately() (ctrl.Result, error) {
	return ctrl.Result{Requeue: true}, nil
}

func NoRequeue() (ctrl.Result, error) {
	return RequeueIfError(nil)
}
