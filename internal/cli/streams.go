// Copyright 2021 NetApp, Inc. All Rights Reserved.

package cli

import (
	"os"

	"github.com/spotinst/ocean-operator/internal/streams"
)

// NewStdIOStreams returns an IOStreams from os.Stdin, os.Stdout.
func NewStdIOStreams() streams.IOStreams {
	return streams.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}
