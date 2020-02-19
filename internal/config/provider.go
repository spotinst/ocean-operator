package config

import (
	"context"
	"fmt"
)

type (
	// Value represents the Operator configuration.
	Value struct {
		// Bootstrap represents the bootstrap configuration.
		Bootstrap ValueBootstrap `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`

		// HTTPClient represents the HTTP client configuration.
		HTTPClient ValueHTTPClient `json:"httpClient,omitempty" yaml:"httpClient,omitempty"`

		// DryRun configures the Operator to print the actions that would be
		// executed, without executing them.
		DryRun bool `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
	}

	// ValueHTTPClient represents the HTTP client configuration.
	ValueHTTPClient struct {
		// BaseURL configures the default base URL of the Spot API.
		BaseURL string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`

		// UserAgent configures the User-Agent HTTP header to set when invoking
		// HTTP requests.
		UserAgent string `json:"userAgent,omitempty" yaml:"userAgent,omitempty"`
	}

	// ValueBootstrap represents the bootstrap configuration.
	ValueBootstrap struct {
		// CRDs is a list of Custom Resource Definitions to install.
		CRDs []ValueCRD `json:"customResourceDefinitions,omitempty" yaml:"customResourceDefinitions,omitempty"`
	}

	// ValueCRD represents a Custom Resource Definition installation configuration.
	ValueCRD struct {
		// Name of the CRD.
		Name string `json:"name,omitempty" yaml:"name,omitempty"`

		// Install policy. One of Always, Never, IfNotPresent.
		InstallPolicy InstallPolicy `json:"installPolicy,omitempty" yaml:"installPolicy,omitempty"`
	}
)

// InstallPolicy is an enumeration of possible installation policies.
type InstallPolicy string

const (
	// InstallAlways indicates always install resource.
	InstallAlways InstallPolicy = "Always"

	// InstallNever indicates never install resource.
	InstallNever InstallPolicy = "Never"

	// InstallIfNotPresent indicates install resource if missing.
	InstallIfNotPresent InstallPolicy = "IfNotPresent"
)

// Provider defines the interface for any component which will provide configuration.
// The Provider should not need to implement its own mutexes, because that will
// be managed by Config.
type Provider interface {
	fmt.Stringer

	// Retrieve retrieves and returns the configuration, or error in case of failure.
	Retrieve(ctx context.Context) (*Value, error)
}
