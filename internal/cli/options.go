// Copyright 2021 NetApp, Inc. All Rights Reserved.

package cli

import (
	"github.com/spotinst/ocean-operator/internal/streams"
	"github.com/spotinst/ocean-operator/pkg/log"
)

// CommonOptions contains common options.
type CommonOptions struct {
	// In, Out, and Err represent the respective data streams that the command
	// may act upon. They are attached directly to any sub-process of the executed
	// command.
	IOStreams streams.IOStreams

	// Log represents the common internal logger implementation.
	Log log.Logger
}

// NewCommonOptions returns a new CommonOptions object.
func NewCommonOptions(streams streams.IOStreams, log log.Logger) *CommonOptions {
	return &CommonOptions{
		IOStreams: streams,
		Log:       log,
	}
}

// SetIOStreams sets the IOStreams.
func (o *CommonOptions) SetIOStreams(streams streams.IOStreams) {
	o.IOStreams = streams
}

// SetLog sets the Log.
func (o *CommonOptions) SetLog(log log.Logger) {
	o.Log = log
}
