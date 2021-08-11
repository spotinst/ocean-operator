// Copyright 2021 NetApp, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"os"

	"github.com/spotinst/ocean-operator/internal/cli"
	operator "github.com/spotinst/ocean-operator/internal/cmd/ocean-operator"
	"github.com/spotinst/ocean-operator/internal/log"
	"github.com/spotinst/ocean-operator/internal/streams"
)

func main() {
	Main()
}

// Main is the actual main(), it invokes Run and exits if an error occurs.
func Main() {
	if err := Run(cli.NewStdIOStreams(), cli.NewZapLogger(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "exited with error: %v\n", err)
		os.Exit(1)
	}
}

// Run invokes the root command and returns an error in case of failure.
func Run(streams streams.IOStreams, log log.Logger, args []string) error {
	cmd := operator.NewCommand(streams, log)
	cmd.SetArgs(args)
	return cmd.Execute()
}
