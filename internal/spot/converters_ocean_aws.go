package spot

import (
	"encoding/json"

	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type oceanAWSClusterConverter struct{}

func (x *oceanAWSClusterConverter) FromObject(cluster *oceanv1.Cluster) (*OceanCluster, error) {
	b, err := json.Marshal(cluster.Spec)
	if err != nil {
		return nil, err
	}

	obj := &aws.Cluster{
		Name:                spotinst.String(cluster.Name),
		ControllerClusterID: spotinst.String(cluster.Name),
	}

	out := &OceanCluster{
		Obj: obj,
	}

	if oceanID := cluster.Status.OceanID; oceanID != "" {
		obj.ID = spotinst.String(oceanID)
		out.ID = oceanID
	}

	if err := json.Unmarshal(b, obj); err != nil {
		return nil, err
	}

	// Manually convert all tags since we have changed the JSON keys.
	if tags := cluster.Spec.Compute.LaunchSpecification.Tags; len(tags) > 0 {
		obj.Compute.LaunchSpecification.Tags = make([]*aws.Tag, len(tags))

		for i, tag := range tags {
			obj.Compute.LaunchSpecification.Tags[i] = &aws.Tag{
				Key:   spotinst.String(tag.Key),
				Value: spotinst.String(tag.Value),
			}
		}
	}

	return out, nil
}

func (x *oceanAWSClusterConverter) ToObject(cluster *OceanCluster) (*oceanv1.Cluster, error) {
	b, err := json.Marshal(cluster.Obj)
	if err != nil {
		return nil, err
	}

	o := &oceanv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: cluster.Name,
		},
		Status: oceanv1.ClusterStatus{
			OceanID: cluster.ID,
		},
	}

	if err := json.Unmarshal(b, &o.Spec); err != nil {
		return nil, err
	}

	// Manually convert all tags since we have changed the JSON keys.
	if tags := cluster.Obj.(*aws.Cluster).Compute.LaunchSpecification.Tags; len(tags) > 0 {
		o.Spec.Compute.LaunchSpecification.Tags = make([]*oceanv1.Tag, len(tags))

		for i, tag := range tags {
			o.Spec.Compute.LaunchSpecification.Tags[i] = &oceanv1.Tag{
				Key:   spotinst.StringValue(tag.Key),
				Value: spotinst.StringValue(tag.Value),
			}
		}
	}

	return o, nil
}

type oceanAWSLaunchSpecConverter struct{}

func (x *oceanAWSLaunchSpecConverter) FromObject(spec *oceanv1.LaunchSpec) (*OceanLaunchSpec, error) {
	b, err := json.Marshal(spec.Spec)
	if err != nil {
		return nil, err
	}

	obj := &aws.LaunchSpec{
		Name: spotinst.String(spec.Name),
	}

	out := &OceanLaunchSpec{
		Obj: obj,
	}

	if specID := spec.Status.LaunchSpecID; specID != "" {
		obj.ID = spotinst.String(specID)
		out.ID = specID
	}

	if err := json.Unmarshal(b, obj); err != nil {
		return nil, err
	}

	return out, nil
}

func (x *oceanAWSLaunchSpecConverter) ToObject(spec *OceanLaunchSpec) (*oceanv1.LaunchSpec, error) {
	b, err := json.Marshal(spec.Obj)
	if err != nil {
		return nil, err
	}

	o := &oceanv1.LaunchSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.Name,
		},
		Status: oceanv1.LaunchSpecStatus{
			LaunchSpecID: spec.ID,
		},
	}

	if err := json.Unmarshal(b, &o.Spec); err != nil {
		return nil, err
	}

	return o, nil
}
