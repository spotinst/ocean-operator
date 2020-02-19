package operator

// OperatorOptions contains options to configure an Operator instance. Each option
// can be set through setter functions. See documentation for each setter function
// for an explanation of the option.
type OperatorOptions struct {
	// HealthzHost and HealthzPort are the host and port the manager should bind
	// to for serving healthz checks.
	HealthzHost string
	HealthzPort int32

	// PrometheusMetricsHost and PrometheusMetricsPort are the host and port the
	// manager should bind to for serving Prometheus metrics.
	PrometheusMetricsHost string
	PrometheusMetricsPort int32

	// OperatorMetricsHost and OperatorMetricsPort are the host and port the
	// manager should bind to for serving custom resource metrics'.
	OperatorMetricsHost string
	OperatorMetricsPort int32
}

// OperatorOption allows specifying various settings configurable by the Operator.
type OperatorOption func(*OperatorOptions)

// WithHealthzHostPort specifies the host and port the manager should bind to
// for serving healthz checks.
func WithHealthzHostPort(host string, port int32) OperatorOption {
	return func(opts *OperatorOptions) {
		opts.HealthzHost = host
		opts.HealthzPort = port
	}
}

// WithPrometheusMetricsHostPort specifies the host and port the manager should
// bind to for serving Prometheus metrics.
func WithPrometheusMetricsHostPort(host string, port int32) OperatorOption {
	return func(opts *OperatorOptions) {
		opts.PrometheusMetricsHost = host
		opts.PrometheusMetricsPort = port
	}
}

// WithOperatorMetricsHostPort specifies the host and port the manager should
// bind to for serving custom resource metrics'.
func WithOperatorMetricsHostPort(host string, port int32) OperatorOption {
	return func(opts *OperatorOptions) {
		opts.OperatorMetricsHost = host
		opts.OperatorMetricsPort = port
	}
}

func defaultOptions() *OperatorOptions {
	return &OperatorOptions{
		HealthzHost:           "0.0.0.0",
		PrometheusMetricsHost: "0.0.0.0",
		OperatorMetricsHost:   "0.0.0.0",
		HealthzPort:           8380,
		PrometheusMetricsPort: 8383,
		OperatorMetricsPort:   8386,
	}
}
