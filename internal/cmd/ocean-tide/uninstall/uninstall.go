// Copyright 2021 NetApp, Inc. All Rights Reserved.

package uninstall

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/pkg/tide"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Options struct {
	*cli.CommonOptions

	ChartNamespace string
	ChartName      string
	Wait           bool
	DryRun         bool
	Timeout        time.Duration

	// internal
	config *rest.Config
}

// NewCommand returns a new cobra.Command for uninstall.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions: commonOptions,
	}

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "UninstallOperator the Ocean Operator and all its components",
		Long:  `UninstallOperator the Ocean Operator and all its components`,
		PreRun: func(cmd *cobra.Command, args []string) {
			cli.PrintFlags(cmd.Flags(), options.Log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run(ctrl.LoggerInto(ctrl.SetupSignalHandler(), options.Log))
		},
	}

	cmd.Flags().StringVar(&options.ChartNamespace, "chart-namespace", oceanv1alpha1.NamespaceSystem, "chart namespace")
	cmd.Flags().StringVar(&options.ChartName, "chart-name", tide.OceanOperatorChart, "chart name")
	cmd.Flags().DurationVar(&options.Timeout, "timeout", 5*time.Minute, "maximum duration before timing out the execution")
	cmd.Flags().BoolVar(&options.Wait, "wait", true, "wait for completion before exiting")
	cmd.Flags().BoolVar(&options.DryRun, "dry-run", false, "only print the actions that would be executed, without executing them")

	return cmd
}

func (x *Options) run(ctx context.Context) (err error) {
	ctrl.SetLogger(x.Log)
	x.config, err = ctrl.GetConfig()
	if err != nil {
		x.Log.Error(err, "unable to get kubeconfig")
		return err
	}

	clientGetter := tide.NewConfigFlags(x.config, x.ChartNamespace)
	manager, err := tide.NewManager(clientGetter, x.Log)
	if err != nil {
		x.Log.Error(err, "unable to create tide manager")
		return err
	}
	deleteOptions := []tide.DeleteOption{
		tide.WithNamespace(x.ChartNamespace),
	}
	if err = manager.DeleteEnvironment(ctx, deleteOptions...); err != nil {
		return err
	}

	operator := tide.NewOperatorOceanComponent(
		tide.WithChartName(x.ChartName),
		tide.WithChartNamespace(x.ChartNamespace),
	)
	if err = tide.UninstallOperator(ctx, operator, clientGetter,
		x.Wait, x.DryRun, x.Timeout, x.Log); err != nil {
		return err
	}

	x.Log.Info("ocean operator is now removed")
	return nil
}
