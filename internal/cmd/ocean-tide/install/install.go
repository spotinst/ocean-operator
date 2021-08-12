// Copyright 2021 NetApp, Inc. All Rights Reserved.

package install

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

	ChartNamespace  string
	ChartName       string
	ChartVersion    string
	ChartURL        string
	ChartValuesJSON string
	Wait            bool
	DryRun          bool
	Timeout         time.Duration

	// internal
	config *rest.Config
}

// NewCommand returns a new cobra.Command for install.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions: commonOptions,
	}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "InstallOperator the Ocean Operator and all its components",
		Long:  `InstallOperator the Ocean Operator and all its components`,
		PreRun: func(cmd *cobra.Command, args []string) {
			cli.PrintFlags(cmd.Flags(), options.Log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run(ctrl.LoggerInto(ctrl.SetupSignalHandler(), options.Log))
		},
	}

	cmd.Flags().StringVar(&options.ChartNamespace, "chart-namespace", oceanv1alpha1.NamespaceSystem, "chart namespace")
	cmd.Flags().StringVar(&options.ChartName, "chart-name", tide.OceanOperatorChart, "chart name")
	cmd.Flags().StringVar(&options.ChartVersion, "chart-version", tide.OceanOperatorVersion, "chart version")
	cmd.Flags().StringVar(&options.ChartURL, "chart-url", tide.OceanOperatorRepository, "chart repository url")
	cmd.Flags().StringVar(&options.ChartValuesJSON, "chart-values", tide.OceanOperatorValues, "chart values (json)")
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
	operator := tide.NewOperatorOceanComponent(
		tide.WithChartName(x.ChartName),
		tide.WithChartNamespace(x.ChartNamespace),
		tide.WithChartURL(x.ChartURL),
		tide.WithChartVersion(x.ChartVersion),
		tide.WithChartValues(x.ChartValuesJSON),
	)
	if err = tide.InstallOperator(ctx, operator, clientGetter,
		x.Wait, x.DryRun, x.Timeout, x.Log); err != nil {
		return err
	}

	x.Log.Info("ocean operator is now installed and managing components")
	return nil
}
