package operator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"runtime"
	"strconv"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	"github.com/spotinst/ocean-operator/internal/config"
	ctrlutil "github.com/spotinst/ocean-operator/internal/util/controller"
	"github.com/spotinst/ocean-operator/internal/version"
	"github.com/spotinst/ocean-operator/pkg/apis"
	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	"github.com/spotinst/ocean-operator/pkg/controller"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	corev1 "k8s.io/api/core/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/intstr"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var log = logf.Log.WithName("operator")

type Operator struct {
	// Name is the Operator name.
	Name string

	// Namespace the Operator runs in.
	OperatorNamespace string

	// Namespaces in which this Operator should manage resources and watch for
	// changes. Accepts multiple comma-separated values. Defaults to all
	// namespaces if empty or unspecified.
	WatchNamespace string

	// Config is the common attributes that can be passed to a Kubernetes client
	// on initialization. Defaults to the configuration loaded from the
	// `KUBECONFIG` environment variable.
	Config *rest.Config

	// Options are additional options to control over the behaviour of the Operator.
	Options []OperatorOption

	// internal
	options *OperatorOptions
	manager manager.Manager
	config  *config.Value
}

// New returns a new Operator instance.
func New(options ...OperatorOption) *Operator {
	return &Operator{Options: options}
}

// Run runs the Operator.
func (x *Operator) Run(ctx context.Context) error {
	steps := []func(ctx context.Context) error{
		x.init,
		x.run,
	}

	for _, fn := range steps {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (x *Operator) run(ctx context.Context) error {
	// Print version information.
	x.printVersion()

	// Load configuration.
	if err := x.loadConfig(ctx); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Install CRDs.
	if err := x.installCRDs(ctx); err != nil {
		return fmt.Errorf("failed to install crds: %w", err)
	}

	// Setup manager.
	if err := x.setupManager(ctx); err != nil {
		return fmt.Errorf("failed to setup manager: %w", err)
	}

	// Setup metrics.
	if err := x.setupMetrics(ctx); err != nil {
		return fmt.Errorf("failed to setup metrics service: %w", err)
	}

	// Start the leadership election and become the leader before proceeding.
	if err := leader.Become(ctx, fmt.Sprintf("%s-lock", x.Name)); err != nil {
		return fmt.Errorf("failed to become the leader: %w", err)
	}

	log.Info("Starting Manager.")

	// Start the manager.
	if err := x.manager.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("manager exited non-zero: %w", err)
	}

	return nil
}

// init initializes the Operator.
func (x *Operator) init(ctx context.Context) error {
	initializers := []func(ctx context.Context) error{
		x.initFlags,
		x.initLogger,
		x.initOptions,
	}

	for _, fn := range initializers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (x *Operator) initFlags(ctx context.Context) error {
	// Add the zap logger flag set to the CLI.
	pflag.CommandLine.AddFlagSet(zap.FlagSet())

	// Add flags registered by imported packages (e.g. glog and controller-runtime).
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Parse flags.
	pflag.Parse()
	return nil
}

func (x *Operator) initLogger(ctx context.Context) error {
	// Use a zap logr.Logger implementation. If none of the zap flags are
	// configured (or if the zap flag set is not being used), this defaults to
	// a production zap logger.
	//
	// The logger instantiated here can be changed to any logger implementing the
	// logr.Logger interface. This logger will be propagated through the whole
	// operator, generating uniform and structured logs.
	logf.SetLogger(zap.Logger())
	return nil
}

func (x *Operator) initOptions(ctx context.Context) error {
	var err error

	// Initialize options.
	x.options = defaultOptions()
	for _, opt := range x.Options {
		opt(x.options)
	}

	// Get the operator name.
	if x.Name == "" {
		x.Name, err = k8sutil.GetOperatorName()
		if err != nil {
			return fmt.Errorf("failed to get operator name: %w", err)
		}
	}

	// Get the namespace the operator is currently deployed in.
	if x.OperatorNamespace == "" {
		x.OperatorNamespace, err = k8sutil.GetOperatorNamespace()
		if err != nil && !errors.Is(err, k8sutil.ErrRunLocal) {
			return fmt.Errorf("failed to get operator namespace: %w", err)
		}
	}

	// Get the namespace the operator should be watching for changes.
	if x.WatchNamespace == "" {
		x.WatchNamespace, err = k8sutil.GetWatchNamespace()
		if err != nil {
			return fmt.Errorf("failed to get watch namespace: %w", err)
		}
	}

	// Load the Kubernetes client config.
	if x.Config == nil {
		x.Config, err = k8sconfig.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load client config: %w", err)
		}
	}

	return nil
}

func (x *Operator) printVersion() {
	log.Info(fmt.Sprintf("Operator Version: %s", version.String()))
	log.Info(fmt.Sprintf("Operator SDK Version: %v", sdkVersion.Version[1:]))
	log.Info(fmt.Sprintf("Spotinst SDK Version: %v", spotinst.SDKVersion))
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()[2:]))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func (x *Operator) loadConfig(ctx context.Context) error {
	log.Info("Loading configuration.")

	apiClient, err := client.New(x.Config, client.Options{})
	if err != nil {
		return fmt.Errorf("failed to instante api client: %w", err)
	}

	x.config, err = ctrlutil.LoadConfig(ctx, apiClient, x.Name, x.OperatorNamespace)
	if err != nil {
		return err
	}

	log.V(1).Info("Configuration loaded", "config", x.config)
	return nil
}

func (x *Operator) installCRDs(ctx context.Context) error {
	log.Info("Installing CRDs.")

	extClient, err := apiextclient.NewForConfig(x.Config)
	if err != nil {
		return fmt.Errorf("failed to instante apiextensions client: %w", err)
	}

	cluster := oceanv1.Cluster{}.CustomResourceDefinition()
	launchSpec := oceanv1.LaunchSpec{}.CustomResourceDefinition()

	for _, obj := range x.config.Bootstrap.CRDs {
		if obj.InstallPolicy == config.InstallNever {
			log.Info("Skipping CRD", "name", obj.Name)
			continue
		}

		log.Info("Installing CRD", "name", obj.Name)
		var crd *apiextv1beta1.CustomResourceDefinition
		switch obj.Name {
		case cluster.Name:
			crd = cluster
		case launchSpec.Name:
			crd = launchSpec
		default:
			return fmt.Errorf("unknown crd name: %v", obj.Name)
		}

		if err := ctrlutil.RegisterCRD(ctx, extClient, crd); err != nil {
			return fmt.Errorf("failed to install crd %s: %w", crd.Name, err)
		}
	}

	return nil
}

func (x *Operator) setupManager(ctx context.Context) (err error) {
	// Create a new manager to provide shared dependencies and start components.
	x.manager, err = manager.New(x.Config, manager.Options{
		Namespace:      x.WatchNamespace,
		MapperProvider: restmapper.NewDynamicRESTMapper,
		HealthProbeBindAddress: net.JoinHostPort(
			x.options.HealthzHost, strconv.Itoa(int(x.options.HealthzPort))),
		MetricsBindAddress: net.JoinHostPort(
			x.options.PrometheusMetricsHost, strconv.Itoa(int(x.options.PrometheusMetricsPort))),
	})
	if err != nil {
		return err
	}

	log.Info("Registering components.")

	// Setup scheme for all resources.
	if err = apis.AddToScheme(x.manager.GetScheme()); err != nil {
		return fmt.Errorf("failed to setup scheme: %w", err)
	}

	// Setup all controllers.
	if err = controller.AddToManager(x.manager); err != nil {
		return fmt.Errorf("failed to setup controllers: %w", err)
	}

	log.Info("Registering healthz checks.")

	// Setup Healthz and Readyz checks.
	if err = x.setupHealthz(ctx); err != nil {
		return fmt.Errorf("failed to setup healthz: %w", err)
	}

	return nil
}

func (x *Operator) setupMetrics(ctx context.Context) error {
	// Serve metrics.
	if err := x.serveCRMetrics(ctx); err != nil {
		log.Info("Could not generate and serve custom resource metrics", "error", err.Error())
	}

	// Metrics ports to expose.
	servicePorts := []corev1.ServicePort{
		{
			Port:     x.options.PrometheusMetricsPort,
			Name:     metrics.OperatorPortName,
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: x.options.PrometheusMetricsPort,
			}},
		{
			Port:     x.options.OperatorMetricsPort,
			Name:     metrics.CRPortName,
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: x.options.OperatorMetricsPort,
			}},
	}

	// Create Service object to expose the metrics port(s).
	service, err := metrics.CreateMetricsService(ctx, x.Config, servicePorts)
	if err != nil {
		log.Info("Could not create metrics service", "error", err.Error())
	}

	// CreateServiceMonitors will automatically create the prometheus-operator
	// ServiceMonitor resources necessary to configure Prometheus to scrape
	// metrics from this operator.
	services := []*corev1.Service{service}
	if _, err = metrics.CreateServiceMonitors(x.Config, x.OperatorNamespace, services); err != nil {
		log.Info("Could not create ServiceMonitor object", "error", err.Error())

		// If this operator is deployed to a cluster without the prometheus-operator
		// running, it will return ErrServiceMonitorNotPresent, which can be used
		// to safely skip ServiceMonitor creation.
		if err == metrics.ErrServiceMonitorNotPresent {
			log.Info("Install prometheus-operator in your cluster to create " +
				"ServiceMonitor objects")
		}
	}

	return nil
}

func (x *Operator) serveCRMetrics(ctx context.Context) error {
	// Get filtered operator/CustomResource specific GVKs.
	filteredGVK, err := k8sutil.GetGVKsFromAddToScheme(apis.AddToScheme)
	if err != nil {
		return err
	}

	// Generate and serve custom resource specific metrics.
	return kubemetrics.GenerateAndServeCRMetrics(x.Config, []string{x.OperatorNamespace},
		filteredGVK, x.options.OperatorMetricsHost, x.options.OperatorMetricsPort)
}

func (x *Operator) setupHealthz(ctx context.Context) error {
	if x.options.HealthzPort <= 0 {
		return nil
	}

	// Setup Healthz check.
	if err := x.manager.AddHealthzCheck("liveness", healthz.Ping); err != nil {
		return err
	}

	// Setup Readyz check.
	if err := x.manager.AddReadyzCheck("readiness", healthz.Ping); err != nil {
		return err
	}

	return nil
}
