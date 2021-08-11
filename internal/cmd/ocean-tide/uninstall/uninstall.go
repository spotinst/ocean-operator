// Copyright 2021 NetApp, Inc. All Rights Reserved.

package uninstall

import (
	"context"

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
		Short: "Uninstall the Ocean Operator and all its components",
		Long:  `Uninstall the Ocean Operator and all its components`,
		PreRun: func(cmd *cobra.Command, args []string) {
			cli.PrintFlags(cmd.Flags(), options.Log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run(ctrl.SetupSignalHandler())
		},
	}

	cmd.Flags().StringVar(&options.ChartNamespace, "chart-namespace", oceanv1alpha1.NamespaceSystem, "chart namespace")
	cmd.Flags().StringVar(&options.ChartName, "chart-name", tide.OceanOperatorChart, "chart name")

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

	operatorOptions := []tide.OperatorChartOption{
		tide.WithOperatorChartName(x.ChartName),
		tide.WithOperatorChartNamespace(x.ChartNamespace),
	}
	operator := tide.NewOperatorOceanComponent(operatorOptions...)
	if err = tide.UninstallOperator(ctx, operator, clientGetter, x.Log); err != nil {
		return err
	}

	x.Log.Info("ocean operator is now removed")
	return nil
}
