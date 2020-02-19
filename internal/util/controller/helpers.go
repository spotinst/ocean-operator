package controller

import (
	"context"
	"os"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/spotinst/ocean-operator/internal/config"
	"github.com/spotinst/ocean-operator/internal/credentials"
	"github.com/spotinst/ocean-operator/internal/spot"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LoadCredentials loads credentials from a chain of providers.
func LoadCredentials(ctx context.Context, client client.Client,
	secretName, namespace string) (*credentials.Value, error) {

	providers := []credentials.Provider{
		&credentials.SecretProvider{
			Client:    client,
			Name:      secretName,
			Namespace: namespace,
		},
		&credentials.SecretProvider{
			Client:    client,
			Name:      secretName,
			Namespace: metav1.NamespaceDefault,
		},
		&credentials.SecretProvider{
			Client:    client,
			Name:      secretName,
			Namespace: metav1.NamespaceSystem,
		},
		&credentials.EnvProvider{},
	}

	return credentials.NewChainProvider(providers...).Get(ctx)
}

// LoadConfig loads configuration from a chain of providers.
func LoadConfig(ctx context.Context, client client.Client,
	configMapName, namespace string) (*config.Value, error) {

	providers := []config.Provider{
		&config.ConfigMapProvider{
			Client:    client,
			Name:      configMapName,
			Namespace: namespace,
		},
		&config.ConfigMapProvider{
			Client:    client,
			Name:      configMapName,
			Namespace: metav1.NamespaceDefault,
		},
		&config.ConfigMapProvider{
			Client:    client,
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
		},
		&config.DefaultProvider{},
	}

	return config.NewChainProvider(providers...).Get(ctx)
}

// NewRequestSpotClient returns a pre-configured Spot client.
func NewRequestSpotClient(ctx RequestContextBase, client client.Client) (spot.Client, error) {
	var (
		clientOpts []spot.ClientOption
		req        = ctx.GetRequest()
		reqLogger  = ctx.GetLogger()
		name       = os.Getenv(k8sutil.OperatorNameEnvVar)
	)

	// Load credentials.
	{
		reqLogger.V(1).Info("Loading credentials")
		creds, err := LoadCredentials(ctx, client, name, req.Namespace)
		if err != nil {
			return nil, err
		}

		clientOpts = append(clientOpts, spot.WithCredentials(
			creds.Token, creds.Account))
	}

	// Load configuration.
	{
		reqLogger.V(1).Info("Loading configuration")
		cfg, err := LoadConfig(ctx, client, name, req.Namespace)
		if err != nil {
			return nil, err
		}

		if dryRun := cfg.DryRun; dryRun {
			clientOpts = append(clientOpts, spot.WithDryRun(dryRun))
		}

		if baseURL := cfg.HTTPClient.BaseURL; len(baseURL) > 0 {
			clientOpts = append(clientOpts, spot.WithBaseURL(baseURL))
		}

		if userAgent := cfg.HTTPClient.UserAgent; len(userAgent) > 0 {
			clientOpts = append(clientOpts, spot.WithUserAgent(userAgent))
		}
	}

	reqLogger.V(1).Info("Instantiating new Spot client")
	return spot.NewClient(clientOpts...), nil
}
