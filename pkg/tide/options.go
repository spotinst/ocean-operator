// Copyright 2021 NetApp, Inc. All Rights Reserved.

package tide

import (
	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
)

// region Options

// ApplyOptions contains apply options.
type ApplyOptions struct {
	ComponentsFilter map[oceanv1alpha1.OceanComponentName]struct{}
}

// DeleteOptions contains delete options.
type DeleteOptions struct{}

// endregion

// region Interfaces

// ApplyOption is some configuration that modifies options for an apply request.
type ApplyOption interface {
	// MutateApplyOptions applies this configuration to the given ApplyOptions.
	MutateApplyOptions(options *ApplyOptions)
}

// DeleteOption is some configuration that modifies options for a delete request.
type DeleteOption interface {
	// MutateDeleteOptions applies this configuration to the given DeleteOptions.
	MutateDeleteOptions(options *DeleteOptions)
}

// endregion

// region Helpers

// ApplyOptionFunc is a convenience type like http.HandlerFunc.
type ApplyOptionFunc func(options *ApplyOptions)

// MutateApplyOptions implements the ApplyOption interface.
func (f ApplyOptionFunc) MutateApplyOptions(options *ApplyOptions) { f(options) }

// DeleteOptionFunc is a convenience type like http.HandlerFunc.
type DeleteOptionFunc func(options *DeleteOptions)

// MutateDeleteOptions implements the DeleteOption interface.
func (f DeleteOptionFunc) MutateDeleteOptions(options *DeleteOptions) { f(options) }

// endregion

// region "Functional" Options

// WithComponentsFilter sets the given ComponentsFilter list.
func WithComponentsFilter(components ...oceanv1alpha1.OceanComponentName) ComponentsFilter {
	return ComponentsFilter{
		components: components,
	}
}

// ComponentsFilter filters components to be applied or deleted.
type ComponentsFilter struct {
	components []oceanv1alpha1.OceanComponentName
}

// MutateApplyOptions implements the ApplyOption interface.
func (w ComponentsFilter) MutateApplyOptions(options *ApplyOptions) {
	options.ComponentsFilter = make(map[oceanv1alpha1.OceanComponentName]struct{})
	for _, component := range w.components {
		options.ComponentsFilter[component] = struct{}{}
	}
}

// Blank assignment to verify that ComponentsFilter implements ApplyOption.
var _ ApplyOption = &ComponentsFilter{}

// endregion
