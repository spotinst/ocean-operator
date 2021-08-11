// Copyright 2021 NetApp, Inc. All Rights Reserved.

package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/internal/version"
)

// NewCommand returns a new cobra.Command for version.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	return &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(commonOptions.IOStreams.Out, version.String())
			return nil
		},
	}
}
