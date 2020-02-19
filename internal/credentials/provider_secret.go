package credentials

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretProvider retrieves credentials from a Secret.
type SecretProvider struct {
	Client          client.Client
	Name, Namespace string
}

// NewSecretProvider returns a new Credentials object wrapping the Secret provider.
func NewSecretProvider(client client.Client, name, namespace string) *Credentials {
	return NewCredentials(&SecretProvider{
		Client:    client,
		Name:      name,
		Namespace: namespace,
	})
}

// Retrieve retrieves and returns the credentials, or error in case of failure.
func (x *SecretProvider) Retrieve(ctx context.Context) (*Value, error) {
	secret, err := getSecret(ctx, x.Client, x.Name, x.Namespace)
	if err != nil {
		return nil, fmt.Errorf("error retrieving secret %q from "+
			"namespace %q: %w", x.Name, x.Namespace, err)
	}

	value, err := decodeSecret(secret)
	if err != nil {
		return nil, fmt.Errorf("error decoding secret %q from "+
			"namespace %q: %w", x.Name, x.Namespace, err)
	}

	return value, nil
}

// String returns the string representation of the Secret provider.
func (x *SecretProvider) String() string {
	return "SecretProvider"
}

func getSecret(ctx context.Context, client client.Client,
	name, namespace string) (*corev1.Secret, error) {

	obj := new(corev1.Secret)
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	return obj, client.Get(ctx, key, obj)
}

func decodeSecret(secret *corev1.Secret) (*Value, error) {
	src := make(map[string]string)
	dst := new(Value)

	// Copy all non-binary secret data.
	for k, v := range secret.StringData {
		src[k] = v
	}

	// Copy all binary secret data.
	for k, v := range secret.Data {
		src[k] = string(v)
	}

	return dst, mapstructure.Decode(src, dst)
}
