// Copyright 2021 NetApp, Inc. All Rights Reserved.

package config

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// DefaultProvider returns the default configuration.
type DefaultProvider struct{}

// NewDefaultProvider returns a new Config object wrapping the Default provider.
func NewDefaultProvider() *Config {
	return NewConfig(&DefaultProvider{})
}

// Retrieve retrieves and returns the configuration, or error in case of failure.
func (x *DefaultProvider) Retrieve(context.Context) (*Value, error) {
	return &Value{ClusterIdentifier: uuid.NewV4().String()}, nil
}

// String returns the string representation of the Default provider.
func (x *DefaultProvider) String() string {
	return "DefaultProvider"
}
