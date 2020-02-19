package spot

import (
	"context"

	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
)

type oceanAWS struct {
	svc aws.Service
}

func (x *oceanAWS) NewClusterConverter() OceanClusterConverter {
	return new(oceanAWSClusterConverter)
}

func (x *oceanAWS) NewLaunchSpecConverter() OceanLaunchSpecConverter {
	return new(oceanAWSLaunchSpecConverter)
}

func (x *oceanAWS) ListClusters(ctx context.Context) ([]*OceanCluster, error) {
	log.V(1).Info("Listing all clusters")

	output, err := x.svc.ListClusters(ctx, &aws.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	clusters := make([]*OceanCluster, len(output.Clusters))
	for i, cluster := range output.Clusters {
		clusters[i] = &OceanCluster{
			TypeMeta: TypeMeta{
				Kind: typeOf(OceanCluster{}),
			},
			ObjectMeta: ObjectMeta{
				ID:        spotinst.StringValue(cluster.ID),
				Name:      spotinst.StringValue(cluster.Name),
				CreatedAt: spotinst.TimeValue(cluster.CreatedAt),
				UpdatedAt: spotinst.TimeValue(cluster.UpdatedAt),
			},
			Obj: cluster,
		}
	}

	return clusters, nil
}

func (x *oceanAWS) ListLaunchSpecs(ctx context.Context) ([]*OceanLaunchSpec, error) {
	log.V(1).Info("Listing all launch specs")

	output, err := x.svc.ListLaunchSpecs(ctx, &aws.ListLaunchSpecsInput{})
	if err != nil {
		return nil, err
	}

	specs := make([]*OceanLaunchSpec, len(output.LaunchSpecs))
	for i, spec := range output.LaunchSpecs {
		specs[i] = &OceanLaunchSpec{
			TypeMeta: TypeMeta{
				Kind: typeOf(OceanLaunchSpec{}),
			},
			ObjectMeta: ObjectMeta{
				ID:        spotinst.StringValue(spec.ID),
				Name:      spotinst.StringValue(spec.Name),
				CreatedAt: spotinst.TimeValue(spec.CreatedAt),
				UpdatedAt: spotinst.TimeValue(spec.UpdatedAt),
			},
			Obj: spec,
		}
	}

	return specs, nil
}

func (x *oceanAWS) GetCluster(ctx context.Context, clusterID string) (*OceanCluster, error) {
	log.V(1).Info("Getting cluster by ID", "clusterID", clusterID)

	input := &aws.ReadClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	output, err := x.svc.ReadCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	cluster := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return cluster, nil
}

func (x *oceanAWS) GetLaunchSpec(ctx context.Context, specID string) (*OceanLaunchSpec, error) {
	log.V(1).Info("Getting launch spec by ID", "specID", specID)

	input := &aws.ReadLaunchSpecInput{
		LaunchSpecID: spotinst.String(specID),
	}

	output, err := x.svc.ReadLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	spec := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanLaunchSpec{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return spec, nil
}

func (x *oceanAWS) CreateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.V(1).Info("Creating a new cluster")

	input := &aws.CreateClusterInput{
		Cluster: cluster.Obj.(*aws.Cluster),
	}

	output, err := x.svc.CreateCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	created := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return created, nil
}

func (x *oceanAWS) CreateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.V(1).Info("Creating a new launch spec")

	input := &aws.CreateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*aws.LaunchSpec),
	}

	output, err := x.svc.CreateLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	created := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return created, nil
}

func (x *oceanAWS) UpdateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.V(1).Info("Updating cluster by ID", "clusterID", cluster.ID)

	input := &aws.UpdateClusterInput{
		Cluster: cluster.Obj.(*aws.Cluster),
	}

	// Remove read-only fields.
	input.Cluster.Region = nil
	input.Cluster.UpdatedAt = nil
	input.Cluster.CreatedAt = nil

	output, err := x.svc.UpdateCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	updated := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return updated, nil
}

func (x *oceanAWS) UpdateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.V(1).Info("Updating launch spec by ID", "specID", spec.ID)

	input := &aws.UpdateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*aws.LaunchSpec),
	}

	// Remove read-only fields.
	input.LaunchSpec.UpdatedAt = nil
	input.LaunchSpec.CreatedAt = nil

	output, err := x.svc.UpdateLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	updated := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return updated, nil
}

func (x *oceanAWS) DeleteCluster(ctx context.Context, clusterID string) error {
	log.V(1).Info("Deleting cluster by ID", "clusterID", clusterID)

	input := &aws.DeleteClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	_, err := x.svc.DeleteCluster(ctx, input)
	return err
}

func (x *oceanAWS) DeleteLaunchSpec(ctx context.Context, specID string) error {
	log.V(1).Info("Deleting launch spec by ID", "specID", specID)

	input := &aws.DeleteLaunchSpecInput{
		LaunchSpecID: spotinst.String(specID),
	}

	_, err := x.svc.DeleteLaunchSpec(ctx, input)
	return err
}
