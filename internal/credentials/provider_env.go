package credentials

import (
	"context"
	"os"

	"github.com/spotinst/ocean-operator/internal/spot"
)

// EnvProvider retrieves credentials from the environment variables of the process.
type EnvProvider struct{}

// NewEnvProvider returns a new Credentials object wrapping the Env provider.
func NewEnvProvider() *Credentials {
	return NewCredentials(&EnvProvider{})
}

// Retrieve retrieves and returns the credentials, or error in case of failure.
func (x *EnvProvider) Retrieve(ctx context.Context) (*Value, error) {
	return &Value{
		Token:   os.Getenv(spot.EnvCredentialsToken),
		Account: os.Getenv(spot.EnvCredentialsAccount),
	}, nil
}

// String returns the string representation of the Env provider.
func (x *EnvProvider) String() string {
	return "EnvProvider"
}
