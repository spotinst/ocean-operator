package spot

import (
	"context"

	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/gcp"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
)

type oceanGCP struct {
	svc gcp.Service
}

func (x *oceanGCP) NewClusterConverter() OceanClusterConverter {
	return new(oceanGCPClusterConverter)
}

func (x *oceanGCP) NewLaunchSpecConverter() OceanLaunchSpecConverter {
	return new(oceanGCPLaunchSpecConverter)
}

func (x *oceanGCP) ListClusters(ctx context.Context) ([]*OceanCluster, error) {
	log.V(1).Info("Listing all clusters")

	output, err := x.svc.ListClusters(ctx, &gcp.ListClustersInput{})
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

func (x *oceanGCP) ListLaunchSpecs(ctx context.Context) ([]*OceanLaunchSpec, error) {
	log.V(1).Info("Listing all launch specs")

	output, err := x.svc.ListLaunchSpecs(ctx, &gcp.ListLaunchSpecsInput{})
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
				ID: spotinst.StringValue(spec.ID),
				//Name:      spotinst.StringValue(spec.Name),
				//CreatedAt: spotinst.TimeValue(spec.CreatedAt),
				//UpdatedAt: spotinst.TimeValue(spec.UpdatedAt),
			},
			Obj: spec,
		}
	}

	return specs, nil
}

func (x *oceanGCP) GetCluster(ctx context.Context, clusterID string) (*OceanCluster, error) {
	log.V(1).Info("Getting cluster by ID", "clusterID", clusterID)

	input := &gcp.ReadClusterInput{
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

func (x *oceanGCP) GetLaunchSpec(ctx context.Context, specID string) (*OceanLaunchSpec, error) {
	log.V(1).Info("Getting launch spec by ID", "specID", specID)

	input := &gcp.ReadLaunchSpecInput{
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
			ID: spotinst.StringValue(output.LaunchSpec.ID),
			//Name:      spotinst.StringValue(output.LaunchSpec.Name),
			//CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			//UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return spec, nil
}

func (x *oceanGCP) CreateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.V(1).Info("Creating a new cluster")

	input := &gcp.CreateClusterInput{
		Cluster: cluster.Obj.(*gcp.Cluster),
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

func (x *oceanGCP) CreateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.V(1).Info("Creating a new launch spec")

	input := &gcp.CreateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*gcp.LaunchSpec),
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
			ID: spotinst.StringValue(output.LaunchSpec.ID),
			//Name:      spotinst.StringValue(output.LaunchSpec.Name),
			//CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			//UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return created, nil
}

func (x *oceanGCP) UpdateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.V(1).Info("Updating cluster by ID", "clusterID", cluster.ID)

	input := &gcp.UpdateClusterInput{
		Cluster: cluster.Obj.(*gcp.Cluster),
	}

	// Remove read-only fields.
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

func (x *oceanGCP) UpdateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.V(1).Info("Updating launch spec by ID", "specID", spec.ID)

	input := &gcp.UpdateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*gcp.LaunchSpec),
	}

	// Remove read-only fields.
	//input.LaunchSpec.UpdatedAt = nil
	//input.LaunchSpec.CreatedAt = nil

	output, err := x.svc.UpdateLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	updated := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID: spotinst.StringValue(output.LaunchSpec.ID),
			//Name:      spotinst.StringValue(output.LaunchSpec.Name),
			//CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			//UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return updated, nil
}

func (x *oceanGCP) DeleteCluster(ctx context.Context, clusterID string) error {
	log.V(1).Info("Deleting cluster by ID", "clusterID", clusterID)

	input := &gcp.DeleteClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	_, err := x.svc.DeleteCluster(ctx, input)
	return err
}

func (x *oceanGCP) DeleteLaunchSpec(ctx context.Context, specID string) error {
	log.V(1).Info("Deleting launch spec by ID", "specID", specID)

	input := &gcp.DeleteLaunchSpecInput{
		LaunchSpecID: spotinst.String(specID),
	}

	_, err := x.svc.DeleteLaunchSpec(ctx, input)
	return err
}
