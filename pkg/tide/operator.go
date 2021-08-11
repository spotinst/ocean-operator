package tide

import (
	"context"
	"fmt"
	"time"

	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/internal/log"
	"github.com/spotinst/ocean-operator/pkg/installer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// InstallOperatorOptions contains install options.
type InstallOperatorOptions struct {
	ClientGetter genericclioptions.RESTClientGetter
	Log          log.Logger
}

// InstallOperator installs the Ocean Operator.
func InstallOperator(
	ctx context.Context,
	operator *oceanv1alpha1.OceanComponent,
	options *InstallOperatorOptions,
) error {
	// install or upgrade
	{
		installerOptions := []installer.InstallerOption{
			installer.WithNamespace(operator.Namespace),
			installer.WithClientGetter(options.ClientGetter),
			installer.WithLogger(options.Log),
		}
		i, err := installer.GetInstance(operator.Spec.Type.String(), installerOptions...)
		if err != nil {
			options.Log.Error(err, "unable to create helm installer")
			return err
		}

		existing, err := i.Get(operator.Spec.Name)
		if err != nil && !installer.IsReleaseNotFound(err) {
			options.Log.Error(err, "error checking ocean operator release")
			return err
		}

		if existing != nil && i.IsUpgrade(operator, existing) {
			options.Log.Info("upgrading ocean operator")
			if _, err = i.Upgrade(operator); err != nil {
				return fmt.Errorf("cannot upgrade ocean operator: %w", err)
			}
		} else {
			options.Log.Info("installing ocean operator")
			if _, err = i.Install(operator); err != nil {
				return fmt.Errorf("cannot install ocean operator: %w", err)
			}
		}
	}

	// validate
	{
		config, err := options.ClientGetter.ToRESTConfig()
		if err != nil {
			return fmt.Errorf("cannot get restconfig: %w", err)
		}

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("cannot connect to cluster: %w", err)
		}

		client := clientSet.AppsV1().Deployments(operator.Namespace)
		err = wait.Poll(5*time.Second, 300*time.Second, func() (bool, error) {
			dep, err := client.Get(ctx, OceanOperatorDeployment, metav1.GetOptions{})
			if err != nil || dep.Status.AvailableReplicas == 0 || dep.Status.UnavailableReplicas != 0 {
				return false, nil
			}
			options.Log.Info("polled",
				"deployment", dep.Name,
				"replicas", dep.Status.AvailableReplicas)
			return true, nil
		})
	}

	return nil
}

// UninstallOperatorOptions contains uninstall options.
type UninstallOperatorOptions struct {
	ClientGetter genericclioptions.RESTClientGetter
	Log          log.Logger
}

// UninstallOperator uninstalls the Ocean Operator.
func UninstallOperator(
	ctx context.Context,
	operator *oceanv1alpha1.OceanComponent,
	options *UninstallOperatorOptions,
) error {
	// uninstall
	{
		installerOptions := []installer.InstallerOption{
			installer.WithNamespace(operator.Namespace),
			installer.WithClientGetter(options.ClientGetter),
			installer.WithLogger(options.Log),
		}
		i, err := installer.GetInstance(operator.Spec.Type.String(), installerOptions...)
		if err != nil {
			options.Log.Error(err, "unable to create helm installer")
			return err
		}

		existing, err := i.Get(operator.Spec.Name)
		if err != nil && !installer.IsReleaseNotFound(err) {
			options.Log.Error(err, "error checking ocean operator release")
			return err
		}

		if existing != nil {
			options.Log.Info("uninstalling ocean operator")
			if err = i.Uninstall(operator); err != nil {
				return fmt.Errorf("cannot uninstall ocean operator: %w", err)
			}
		}
	}

	// validate
	{
		config, err := options.ClientGetter.ToRESTConfig()
		if err != nil {
			return fmt.Errorf("cannot get restconfig: %w", err)
		}

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("cannot connect to cluster: %w", err)
		}

		client := clientSet.AppsV1().Deployments(operator.Namespace)
		err = wait.Poll(5*time.Second, 300*time.Second, func() (bool, error) {
			dep, err := client.Get(ctx, OceanOperatorDeployment, metav1.GetOptions{})
			if err == nil {
				options.Log.Info("polled",
					"deployment", dep.Name,
					"replicas", dep.Status.AvailableReplicas)
				return false, nil
			}
			return true, nil
		})
	}

	return nil
}

const (
	OceanOperatorDeployment = "ocean-operator"
	OceanOperatorConfigMap  = "ocean-operator"
	OceanOperatorSecret     = "ocean-operator"
	OceanOperatorChart      = "ocean-operator"
	OceanOperatorRepository = "https://charts.spot.io"
	OceanOperatorVersion    = "" // empty string indicates the latest chart version
	OceanOperatorValues     = ""

	LegacyOceanControllerDeployment = "spotinst-kubernetes-cluster-controller"
	LegacyOceanControllerSecret     = "spotinst-kubernetes-cluster-controller"
	LegacyOceanControllerConfigMap  = "spotinst-kubernetes-cluster-controller-config"
)

// NewOperatorOceanComponent returns an oceanv1alpha1.OceanComponent
// representing the Ocean Operator.
func NewOperatorOceanComponent(options ...OperatorChartOption) *oceanv1alpha1.OceanComponent {
	comp := &oceanv1alpha1.OceanComponent{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: oceanv1alpha1.NamespaceSystem,
			Name:      OceanOperatorChart,
		},
		Spec: oceanv1alpha1.OceanComponentSpec{
			Type:    oceanv1alpha1.OceanComponentTypeHelm,
			State:   oceanv1alpha1.OceanComponentStatePresent,
			Name:    OceanOperatorChart,
			URL:     OceanOperatorRepository,
			Version: OceanOperatorVersion,
		},
	}

	opts := new(OperatorChartOptions)
	for _, opt := range options {
		opt.MutateOperatorChartOptions(opts)
	}

	comp.Spec.Name = oceanv1alpha1.OceanComponentName(opts.Name)
	comp.Namespace = opts.Namespace
	comp.Spec.URL = opts.URL
	comp.Spec.Version = opts.Version
	comp.Spec.Values = opts.Values

	return comp
}

// region Options

// OperatorChartOptions contains Ocean Operator chart options.
type OperatorChartOptions struct {
	Name      string
	Namespace string
	URL       string
	Version   string
	Values    string
}

// endregion

// region Interfaces

// OperatorChartOption is some configuration that modifies options for an Operator.
type OperatorChartOption interface {
	// MutateOperatorChartOptions applies this configuration to the given OperatorChartOptions.
	MutateOperatorChartOptions(options *OperatorChartOptions)
}

// endregion

// region Helpers

// OperatorChartOptionFunc is a convenience type like http.HandlerFunc.
type OperatorChartOptionFunc func(options *OperatorChartOptions)

// MutateOperatorChartOptions implements the OperatorChartOption interface.
func (f OperatorChartOptionFunc) MutateOperatorChartOptions(options *OperatorChartOptions) {
	f(options)
}

// endregion

// region "Functional" Options

// WithOperatorChartName sets the given chart name.
func WithOperatorChartName(name string) OperatorChartOption {
	return OperatorChartOptionFunc(func(options *OperatorChartOptions) {
		options.Name = name
	})
}

// WithOperatorChartNamespace sets the given chart namespace.
func WithOperatorChartNamespace(namespace string) OperatorChartOption {
	return OperatorChartOptionFunc(func(options *OperatorChartOptions) {
		options.Namespace = namespace
	})
}

// WithOperatorChartURL sets the given chart URL.
func WithOperatorChartURL(url string) OperatorChartOption {
	return OperatorChartOptionFunc(func(options *OperatorChartOptions) {
		options.URL = url
	})
}

// WithOperatorChartVersion sets the given chart version.
func WithOperatorChartVersion(version string) OperatorChartOption {
	return OperatorChartOptionFunc(func(options *OperatorChartOptions) {
		options.Version = version
	})
}

// WithOperatorChartValues sets the given chart values.
func WithOperatorChartValues(values string) OperatorChartOption {
	return OperatorChartOptionFunc(func(options *OperatorChartOptions) {
		options.Values = values
	})
}

// endregion
