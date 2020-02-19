package spot

import (
	"context"
	"errors"

	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	"github.com/spotinst/spotinst-sdk-go/spotinst/credentials"
)

// ErrNotImplemented is the error returned if a method is not implemented.
var ErrNotImplemented = errors.New("spotinst: not implemented")

const (
	// EnvCredentialsToken specifies the name of the environment variable points
	// to the Spot Token.
	EnvCredentialsToken = credentials.EnvCredentialsVarToken

	// EnvCredentialsAccount specifies the name of the environment variable points
	// to the Spot account ID.
	EnvCredentialsAccount = credentials.EnvCredentialsVarAccount

	// EnvBaseURL specifies the name of the environment variable points
	// to the Spot API base URL.
	EnvBaseURL = "SPOTINST_BASE_URL"
)

type (
	// Client defines the interface of the Spot API.
	Client interface {
		// Accounts returns an instance of Accounts interface.
		Accounts() Accounts

		// Services returns an instance of Services interface.
		Services() Services
	}

	// Accounts defines the interface of the Spot Accounts API.
	Accounts interface {
		// ListAccounts returns a list of Spot accounts.
		ListAccounts(ctx context.Context) ([]*Account, error)
	}

	// Services defines the interface of the Spot Services API.
	Services interface {
		// Ocean returns an instance of Ocean interface by cloud provider and
		// orchestrator names.
		Ocean(provider CloudProviderName) (Ocean, error)
	}

	// Ocean defines the interface of the Spot Ocean API.
	Ocean interface {
		// NewClusterConverter returns new instance of OceanClusterConverter
		// interface for converting Cluster objects.
		NewClusterConverter() OceanClusterConverter

		// NewLaunchSpecConverter returns new instance of OceanLaunchSpecConverter
		// interface for converting LaunchSpec objects.
		NewLaunchSpecConverter() OceanLaunchSpecConverter

		// ListClusters returns a list of Ocean clusters.
		ListClusters(ctx context.Context) ([]*OceanCluster, error)

		// ListLaunchSpecs returns a list of Ocean launch specs.
		ListLaunchSpecs(ctx context.Context) ([]*OceanLaunchSpec, error)

		// GetCluster returns an Ocean cluster spec by ID.
		GetCluster(ctx context.Context, clusterID string) (*OceanCluster, error)

		// GetLaunchSpec returns an Ocean launch spec by ID.
		GetLaunchSpec(ctx context.Context, specID string) (*OceanLaunchSpec, error)

		// CreateCluster creates a new Ocean cluster.
		CreateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error)

		// CreateLaunchSpec creates a new Ocean launch spec.
		CreateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error)

		// UpdateCluster updates an existing Ocean cluster by ID.
		UpdateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error)

		// UpdateLaunchSpec updates an existing Ocean launch spec by ID.
		UpdateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error)

		// DeleteCluster deletes an Ocean cluster spec by ID.
		DeleteCluster(ctx context.Context, clusterID string) error

		// DeleteLaunchSpec deletes an Ocean launch spec by ID.
		DeleteLaunchSpec(ctx context.Context, specID string) error
	}

	// OceanClusterConverter is the interface that every Ocean cluster concrete
	// implementation should obey.
	OceanClusterConverter interface {
		// FromObject converts from Cluster to OceanCluster.
		FromObject(*oceanv1.Cluster) (*OceanCluster, error)

		// ToObject converts from OceanCluster to Cluster.
		ToObject(*OceanCluster) (*oceanv1.Cluster, error)
	}

	// OceanLaunchSpecConverter is the interface that every Ocean launch spec
	// concrete implementation should obey.
	OceanLaunchSpecConverter interface {
		// FromObject converts from LaunchSpec to OceanLaunchSpec.
		FromObject(*oceanv1.LaunchSpec) (*OceanLaunchSpec, error)

		// ToObject converts from OceanLaunchSpec to LaunchSpec.
		ToObject(*OceanLaunchSpec) (*oceanv1.LaunchSpec, error)
	}

	// CloudProviderName represents the name of a cloud provider.
	CloudProviderName string
)

// Cloud Providers.
const (
	CloudProviderAWS   CloudProviderName = "aws"
	CloudProviderGCP   CloudProviderName = "gcp"
	CloudProviderAzure CloudProviderName = "azure"
)
