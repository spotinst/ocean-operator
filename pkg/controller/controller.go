package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager.
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager.
func AddToManager(manager manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(manager); err != nil {
			return err
		}
	}
	return nil
}
