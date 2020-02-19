package config

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigMapProvider retrieves configuration from a ConfigMap.
type ConfigMapProvider struct {
	Client          client.Client
	Name, Namespace string
}

// NewConfigMapProvider returns a new Config object wrapping the ConfigMap provider.
func NewConfigMapProvider(client client.Client, name, namespace string) *Config {
	return NewConfig(&ConfigMapProvider{
		Client:    client,
		Name:      name,
		Namespace: namespace,
	})
}

// Retrieve retrieves and returns the configuration, or error in case of failure.
func (x *ConfigMapProvider) Retrieve(ctx context.Context) (*Value, error) {
	configMap, err := getConfigMap(ctx, x.Client, x.Name, x.Namespace)
	if err != nil {
		return nil, fmt.Errorf("error retrieving configmap %q from "+
			"namespace %q: %w", x.Name, x.Namespace, err)
	}

	value, err := decodeConfigMap(configMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding configmap %q from "+
			"namespace %q: %w", x.Name, x.Namespace, err)
	}

	return value, nil
}

// String returns the string representation of the ConfigMap provider.
func (x *ConfigMapProvider) String() string {
	return "ConfigMapProvider"
}

func getConfigMap(ctx context.Context, client client.Client,
	name, namespace string) (*corev1.ConfigMap, error) {

	obj := new(corev1.ConfigMap)
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	return obj, client.Get(ctx, key, obj)
}

func decodeConfigMap(configMap *corev1.ConfigMap) (*Value, error) {
	src := []byte(configMap.Data["config.yaml"])
	dst := new(Value)

	return dst, yaml.Unmarshal(src, dst)
}
