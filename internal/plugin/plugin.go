package plugin

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spotinst/ocean-operator/internal/version"
	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl-ocean",
		Short:         "kubectl plugin for managing Ocean resources.",
		Long:          `kubectl plugin for managing Ocean resources.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmds := []func(streams *genericclioptions.IOStreams) *cobra.Command{
		NewGetCmd,
		NewCreateCmd,
		NewUpdateCmd,
		NewDeleteCmd,
		NewVersionCmd,
	}

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd(streams))
	}

	return cmd
}

func NewVersionCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Print version information",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(streams.Out, version.String())
			return err
		},
		Example: `
# Print version information of the plugin.
kubectl ocean version`,
	}

	return cmd
}

func NewCreateCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create -f FILENAME",
		Short:        "Create a resource from a file",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := newExecCmd(streams, "kubectl", "create")

			if file, _ := c.Flags().GetString("file"); file != "" {
				cmd.Args = append(cmd.Args, "-f", file)
			}

			return cmd.Run()
		},
		Example: `
# Create a cluster using the type and name specified in manifest file.
kubectl ocean create -f cluster.yaml

# Create a launchspec using the type and name specified in manifest file.
kubectl ocean create -f spec.yaml`,
	}

	addFlags(cmd.Flags())
	return cmd
}

func NewUpdateCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update -f FILENAME",
		Short:        "Update a resource from a file",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := newExecCmd(streams, "kubectl", "apply")

			if file, _ := c.Flags().GetString("file"); file != "" {
				cmd.Args = append(cmd.Args, "-f", file)
			}

			return cmd.Run()
		},
		Example: `
# Update a cluster using the type and name specified in manifest file.
kubectl ocean update -f cluster.yaml

# Update a launchspec using the type and name specified in manifest file.
kubectl ocean update -f spec.yaml`,
	}

	addFlags(cmd.Flags())
	return cmd
}

func NewDeleteCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete",
		Short:        "Delete resources by filenames and names",
		SilenceUsage: true,
	}

	cmds := []func(streams *genericclioptions.IOStreams) *cobra.Command{
		newDeleteClusterCmd,
		newDeleteLaunchSpecCmd,
	}

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd(streams))
	}

	return cmd
}

func NewGetCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Display one or many resources",
		SilenceUsage: true,
	}

	cmds := []func(streams *genericclioptions.IOStreams) *cobra.Command{
		newGetClusterCmd,
		newGetLaunchSpecCmd,
	}

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd(streams))
	}

	return cmd
}

func newGetClusterCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:          oceanv1.ClusterSingularName,
		Aliases:      []string{oceanv1.ClusterPluralName, "c", "cs"},
		Short:        "Display one or many clusters",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return newExecCmd(streams, "kubectl", "get",
				oceanv1.Cluster{}.CustomResourceDefinition().Name).Run()
		},
		Example: `
# Get a cluster using cluster name.
kubectl ocean get cluster cluster01`,
	}
}

func newGetLaunchSpecCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:          oceanv1.LaunchSpecSingularName,
		Aliases:      []string{oceanv1.LaunchSpecPluralName, "specs", "spec", "l", "ls"},
		Short:        "Display one or many launchspecs",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return newExecCmd(streams, "kubectl", "get",
				oceanv1.LaunchSpec{}.CustomResourceDefinition().Name).Run()
		},
		Example: `
# Get a launchspec using spec name.
kubectl ocean get launchspec spec01`,
	}
}

func newDeleteClusterCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          fmt.Sprintf("%s (-f FILENAME | NAME)", oceanv1.ClusterSingularName),
		Aliases:      []string{oceanv1.ClusterPluralName, "c", "cs"},
		Short:        "Delete clusters by filenames and names",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := newExecCmd(streams, "kubectl", "delete")

			if file, _ := c.Flags().GetString("file"); file != "" {
				cmd.Args = append(cmd.Args, "-f", file)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("cluster name can't be empty")
				}

				cmd.Args = append(cmd.Args,
					oceanv1.Cluster{}.CustomResourceDefinition().Name, args[0])
			}

			return cmd.Run()
		},
		Example: `
# Delete a cluster using the type and name specified in manifest file.
kubectl ocean delete cluster -f cluster.yaml

# Delete a cluster using cluster name.
kubectl ocean delete cluster cluster01`,
	}

	addFlags(cmd.Flags())
	return cmd
}

func newDeleteLaunchSpecCmd(streams *genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          fmt.Sprintf("%s (-f FILENAME | NAME)", oceanv1.LaunchSpecSingularName),
		Aliases:      []string{oceanv1.LaunchSpecPluralName, "specs", "spec", "l", "ls"},
		Short:        "Delete launchspecs by filenames and names",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := newExecCmd(streams, "kubectl", "delete")

			if file, _ := c.Flags().GetString("file"); file != "" {
				cmd.Args = append(cmd.Args, "-f", file)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("launchspec name can't be empty")
				}

				cmd.Args = append(cmd.Args,
					oceanv1.LaunchSpec{}.CustomResourceDefinition().Name, args[0])
			}

			return cmd.Run()
		},
		Example: `
# Delete a launchspec using the type and name specified in manifest file.
kubectl ocean delete launchspec -f spec.yaml

# Delete a launchspec using spec name.
kubectl ocean delete launchspec spec01`,
	}

	addFlags(cmd.Flags())
	return cmd
}

func newExecCmd(streams *genericclioptions.IOStreams, name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	// IO streams.
	cmd.Stdout = streams.Out
	cmd.Stderr = streams.ErrOut
	cmd.Stdin = streams.In

	return cmd
}

func addFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP("file", "f", "", "path to the resource file")
}
