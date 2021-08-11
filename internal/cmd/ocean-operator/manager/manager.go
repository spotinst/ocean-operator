// Copyright 2021 NetApp, Inc. All Rights Reserved.

package manager

import (
	"context"
	"fmt"
	goruntime "runtime"

	"github.com/spf13/cobra"
	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/controllers"
	"github.com/spotinst/ocean-operator/internal/cli"
	"github.com/spotinst/ocean-operator/internal/ocean"
	"github.com/spotinst/ocean-operator/internal/version"
	"github.com/spotinst/ocean-operator/pkg/tide"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	//+kubebuilder:scaffold:imports
)

type Options struct {
	*cli.CommonOptions

	Namespace      string
	MetricsAddress string
	ProbeAddress   string
	LeaderLock     string
	LeaderElection bool
	Components     *ocean.ComponentsFlag

	// internal
	config  *rest.Config
	manager manager.Manager
}

// NewCommand returns a new cobra.Command for manager.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions: commonOptions,
		Components:    ocean.NewDefaultComponentsFlag(commonOptions.Log),
	}

	cmd := &cobra.Command{
		Use:           "manager",
		Short:         "Start the controller runtime manager",
		Long:          `Start the controller runtime manager`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			cli.PrintFlags(cmd.Flags(), options.Log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run(ctrl.SetupSignalHandler())
		},
	}

	cmd.Flags().StringVar(&options.Namespace, "namespace", oceanv1alpha1.NamespaceSystem, "namespace the operator runs in")
	cmd.Flags().StringVar(&options.MetricsAddress, "metrics-bind-address", ":8080", "address the metric endpoint binds to")
	cmd.Flags().StringVar(&options.ProbeAddress, "health-probe-bind-address", ":8081", "address the probe endpoint binds to")
	cmd.Flags().StringVar(&options.LeaderLock, "leader-lock", "6c511c84.spot.io", "leader election lock name")
	cmd.Flags().BoolVar(&options.LeaderElection, "leader-elect", false, "enable leader election")
	cmd.Flags().Var(options.Components, "components", "list of components to manage")

	return cmd
}

func (x *Options) run(ctx context.Context) error {
	for _, fn := range []func(context.Context) error{
		x.printVersion,
		x.setupConfig,
		x.setupOperator,
		x.setupManager,
		x.setupChecks,
		x.startManager,
	} {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (x *Options) printVersion(ctx context.Context) error {
	x.Log.Info(fmt.Sprintf("operator version: %s", version.String()))
	x.Log.Info(fmt.Sprintf("go version: %s", goruntime.Version()[2:]))
	x.Log.Info(fmt.Sprintf("go os/arch: %s/%s", goruntime.GOOS, goruntime.GOARCH))
	return nil
}

func (x *Options) setupConfig(ctx context.Context) error {
	ctrl.SetLogger(x.Log)

	var err error
	x.config, err = ctrl.GetConfig()
	if err != nil {
		x.Log.Error(err, "unable to get kubeconfig")
		return err
	}

	return nil
}

func (x *Options) setupOperator(ctx context.Context) error {
	clientGetter := tide.NewConfigFlags(x.config, x.Namespace)
	manager, err := tide.NewManager(clientGetter, x.Log)
	if err != nil {
		x.Log.Error(err, "unable to create tide manager")
		return err
	}

	applyOptions := []tide.ApplyOption{
		tide.WithComponentsFilter(x.Components.List()...),
	}
	if err = manager.ApplyEnvironment(ctx, applyOptions...); err != nil {
		return err
	}

	return nil
}

func (x *Options) setupManager(ctx context.Context) error {
	var err error
	x.manager, err = ctrl.NewManager(x.config, ctrl.Options{
		Scheme:                 tide.DefaultScheme(),
		Port:                   9443,
		MetricsBindAddress:     x.MetricsAddress,
		HealthProbeBindAddress: x.ProbeAddress,
		LeaderElection:         x.LeaderElection,
		LeaderElectionID:       x.LeaderLock,
	})
	if err != nil {
		x.Log.Error(err, "unable to create runtime manager")
		return err
	}

	if err = (&controllers.OceanComponentReconciler{
		Scheme:       x.manager.GetScheme(),
		Client:       x.manager.GetClient(),
		ClientGetter: tide.NewConfigFlags(x.config, x.Namespace),
		Log:          x.Log.WithName("oceancomponent"),
	}).SetupWithManager(x.manager); err != nil {
		x.Log.Error(err, "unable to create controller", "controller", "oceancomponent")
		return err
	}

	//+kubebuilder:scaffold:builder
	return nil
}

func (x *Options) setupChecks(ctx context.Context) error {
	x.Log.Info("registering checks")
	if err := x.manager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		x.Log.Error(err, "unable to set up health check")
		return err
	}
	if err := x.manager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		x.Log.Error(err, "unable to set up ready check")
		return err
	}
	return nil
}

func (x *Options) startManager(ctx context.Context) error {
	x.Log.Info("starting manager")
	if err := x.manager.Start(ctx); err != nil {
		x.Log.Error(err, "manager exited non-zero")
		return err
	}
	return nil
}
