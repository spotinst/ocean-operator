// Copyright 2021 NetApp, Inc. All Rights Reserved.

package install

import (
	"context"

	"github.com/spf13/cobra"
	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/internal/ocean"
	"github.com/spotinst/ocean-operator/pkg/tide"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Options struct {
	*cli.CommonOptions

	Namespace       string
	ChartName       string
	ChartVersion    string
	ChartURL        string
	ChartValuesJSON string
	Components      *ocean.ComponentsFlag

	// internal
	config *rest.Config
}

// NewCommand returns a new cobra.Command for install.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions: commonOptions,
		Components:    ocean.NewDefaultComponentsFlag(commonOptions.Log),
	}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the Ocean Operator and all its components",
		Long:  `Install the Ocean Operator and all its components`,
		PreRun: func(cmd *cobra.Command, args []string) {
			cli.PrintFlags(cmd.Flags(), options.Log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run(ctrl.SetupSignalHandler())
		},
	}

	cmd.Flags().StringVar(&options.Namespace, "namespace", oceanv1alpha1.NamespaceSystem, "namespace the operator should be installed in")
	cmd.Flags().StringVar(&options.ChartName, "chart-name", tide.OceanOperatorChart, "chart name")
	cmd.Flags().StringVar(&options.ChartVersion, "chart-version", tide.OceanOperatorVersion, "chart version")
	cmd.Flags().StringVar(&options.ChartURL, "chart-url", tide.OceanOperatorRepository, "chart repository url")
	cmd.Flags().StringVar(&options.ChartValuesJSON, "chart-values", tide.OceanOperatorValues, "chart values (json)")
	cmd.Flags().Var(options.Components, "components", "list of components to install")

	return cmd
}

func (x *Options) run(ctx context.Context) error {
	ctrl.SetLogger(x.Log)

	var err error
	x.config, err = ctrl.GetConfig()
	if err != nil {
		x.Log.Error(err, "unable to get kubeconfig")
		return err
	}

	clientGetter := tide.NewConfigFlags(x.config, x.Namespace)
	manager, err := tide.NewManager(clientGetter, x.Log)
	if err != nil {
		x.Log.Error(err, "unable to create tide manager")
		return err
	}

	operatorOptions := []tide.OperatorChartOption{
		tide.WithOperatorChartName(x.ChartName),
		tide.WithOperatorChartNamespace(x.Namespace),
		tide.WithOperatorChartURL(x.ChartURL),
		tide.WithOperatorChartVersion(x.ChartVersion),
		tide.WithOperatorChartValues(x.ChartValuesJSON),
	}
	operator := tide.NewOperatorOceanComponent(operatorOptions...)
	installOptions := &tide.InstallOperatorOptions{
		ClientGetter: clientGetter,
		Log:          x.Log,
	}
	if err = tide.InstallOperator(ctx, operator, installOptions); err != nil {
		return err
	}

	applyOptions := []tide.ApplyOption{
		tide.WithComponentsFilter(x.Components.List()...),
	}
	if err = manager.ApplyEnvironment(ctx, applyOptions...); err != nil {
		return err
	}

	x.Log.Info("ocean operator is now installed and managing components")
	return nil
}
