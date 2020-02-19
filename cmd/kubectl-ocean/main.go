package main

import (
	"os"

	"github.com/spotinst/ocean-operator/internal/plugin"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	streams := &genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	if err := plugin.NewCmd(streams).Execute(); err != nil {
		os.Exit(1)
	}
}
