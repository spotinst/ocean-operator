package credentials

import (
	"context"
	"fmt"
)

// Value represents the Operator credentials.
type Value struct {
	Token   string `json:"token,omitempty"   yaml:"token,omitempty"`
	Account string `json:"account,omitempty" yaml:"account,omitempty"`
}

// Provider defines the interface for any component which will provide credentials.
// The Provider should not need to implement its own mutexes, because that will
// be managed by Config.
type Provider interface {
	fmt.Stringer

	// Retrieve retrieves and returns the credentials, or error in case of failure.
	Retrieve(ctx context.Context) (*Value, error)
}
