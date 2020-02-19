package config

import (
	"context"
	"fmt"

	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
)

// DefaultProvider returns the default configuration.
type DefaultProvider struct{}

// NewDefaultProvider returns a new Config object wrapping the Default provider.
func NewDefaultProvider() *Config {
	return NewConfig(&DefaultProvider{})
}

// Retrieve retrieves and returns the configuration, or error in case of failure.
func (x *DefaultProvider) Retrieve(context.Context) (*Value, error) {
	return &Value{
		Bootstrap: ValueBootstrap{
			CRDs: []ValueCRD{
				{
					Name: fmt.Sprintf("%s.%s",
						oceanv1.ClusterPluralName,
						oceanv1.SchemeGroupVersion.Group),
					InstallPolicy: InstallAlways,
				},
				{
					Name: fmt.Sprintf("%s.%s",
						oceanv1.LaunchSpecPluralName,
						oceanv1.SchemeGroupVersion.Group),
					InstallPolicy: InstallAlways,
				},
			},
		},
	}, nil
}

// String returns the string representation of the Default provider.
func (x *DefaultProvider) String() string {
	return "DefaultProvider"
}
