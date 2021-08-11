// Copyright 2021 NetApp, Inc. All Rights Reserved.

package operator

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/internal/cmd/ocean-operator/manager"
	"github.com/spotinst/ocean-operator/internal/cmd/ocean-operator/version"
	"github.com/spotinst/ocean-operator/internal/log"
	"github.com/spotinst/ocean-operator/internal/streams"
)

func NewCommand(streams streams.IOStreams, log log.Logger) *cobra.Command {
	// Initialize common options.
	options := cli.NewCommonOptions(streams, log)

	// Root command.
	cmd := &cobra.Command{
		Use:          "ocean-operator",
		Short:        `Ocean Operator manager`,
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			options.SetLog(cli.NewZapLogger())
		},
	}

	// Sub-commands.
	cmd.AddCommand(version.NewCommand(options))
	cmd.AddCommand(manager.NewCommand(options))

	// IO streams.
	cmd.SetIn(streams.In)
	cmd.SetOut(streams.Out)
	cmd.SetErr(streams.ErrOut)

	// Add flags registered by imported packages (e.g. klog and controller-runtime).
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}
