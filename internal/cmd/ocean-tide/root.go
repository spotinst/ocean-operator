// Copyright 2021 NetApp, Inc. All Rights Reserved.

package tide

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/internal/cmd/ocean-tide/install"
	"github.com/spotinst/ocean-operator/internal/cmd/ocean-tide/uninstall"
	"github.com/spotinst/ocean-operator/internal/cmd/ocean-tide/version"
	"github.com/spotinst/ocean-operator/internal/streams"
	"github.com/spotinst/ocean-operator/pkg/log"
)

func NewCommand(streams streams.IOStreams, log log.Logger) *cobra.Command {
	// Initialize common options.
	options := cli.NewCommonOptions(streams, log)

	// Root command.
	cmd := &cobra.Command{
		Use:          "ocean-tide",
		Short:        `Ocean Operator installer`,
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
	cmd.AddCommand(install.NewCommand(options))
	cmd.AddCommand(uninstall.NewCommand(options))

	// IO streams.
	cmd.SetIn(streams.In)
	cmd.SetOut(streams.Out)
	cmd.SetErr(streams.ErrOut)

	// Add flags registered by imported packages (e.g. klog and controller-runtime).
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}
