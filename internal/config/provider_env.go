// Copyright 2021 NetApp, Inc. All Rights Reserved.

package config

import (
	"context"
	"os"
)

// EnvClusterIdentifier specifies the name of the environment variable points
// to cluster identifier.
const EnvClusterIdentifier = "SPOTINST_CLUSTER_IDENTIFIER"

// EnvProvider retrieves configuration from the environment variables of the process.
type EnvProvider struct{}

// NewEnvProvider returns a new Config object wrapping the Env provider.
func NewEnvProvider() *Config {
	return NewConfig(&EnvProvider{})
}

// Retrieve retrieves and returns the configuration, or error in case of failure.
func (x *EnvProvider) Retrieve(ctx context.Context) (*Value, error) {
	return &Value{ClusterIdentifier: os.Getenv(EnvClusterIdentifier)}, nil
}

// String returns the string representation of the Env provider.
func (x *EnvProvider) String() string {
	return "EnvProvider"
}
