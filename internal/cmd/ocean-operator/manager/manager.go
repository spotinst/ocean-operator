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

	LeaderElection      bool
	LeaderLock          string
	MetricsAddress      string
	ProbeAddress        string
	BootstrapNamespace  string
	BootstrapComponents *ocean.ComponentsFlag

	// internal
	config  *rest.Config
	manager manager.Manager
}

// NewCommand returns a new cobra.Command for manager.
func NewCommand(commonOptions *cli.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions:       commonOptions,
		BootstrapComponents: ocean.NewEmptyComponentsFlag(commonOptions.Log),
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

	// bind
	cmd.Flags().StringVar(&options.MetricsAddress, "metrics-bind-address", ":8080", "address the metric endpoint binds to")
	cmd.Flags().StringVar(&options.ProbeAddress, "health-probe-bind-address", ":8081", "address the probe endpoint binds to")

	// leadership
	cmd.Flags().BoolVar(&options.LeaderElection, "leader-elect", false, "enable leader election")
	cmd.Flags().StringVar(&options.LeaderLock, "leader-lock", "6c511c84.spot.io", "leader election lock name")

	// bootstrap
	cmd.Flags().StringVar(&options.BootstrapNamespace, "bootstrap-namespace", oceanv1alpha1.NamespaceSystem, "namespace where components should be installed during environment bootstrapping")
	cmd.Flags().Var(options.BootstrapComponents, "bootstrap-components", "list of components to install during environment bootstrapping")

	return cmd
}

func (x *Options) run(ctx context.Context) error {
	for _, fn := range []func(context.Context) error{
		x.printVersion,
		x.setupConfig,
		x.setupEnvironment,
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

func (x *Options) setupConfig(ctx context.Context) (err error) {
	ctrl.SetLogger(x.Log)
	x.config, err = ctrl.GetConfig()
	if err != nil {
		x.Log.Error(err, "unable to get kubeconfig")
		return err
	}
	return nil
}

func (x *Options) setupEnvironment(ctx context.Context) error {
	clientGetter := tide.NewConfigFlags(x.config, x.BootstrapNamespace)
	manager, err := tide.NewManager(clientGetter, x.Log)
	if err != nil {
		x.Log.Error(err, "unable to create tide manager")
		return err
	}

	applyOptions := []tide.ApplyOption{
		tide.WithNamespace(x.BootstrapNamespace),
		tide.WithComponentsFilter(x.BootstrapComponents.List()...),
	}
	if err = manager.ApplyEnvironment(ctx, applyOptions...); err != nil {
		return err
	}

	return nil
}

func (x *Options) setupManager(ctx context.Context) (err error) {
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
		ClientGetter: tide.NewConfigFlags(x.config, x.BootstrapNamespace),
		Log:          x.Log.WithName("oceancomponent"),
		Namespace:    x.BootstrapNamespace,
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
