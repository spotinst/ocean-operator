package spot

import (
	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
)

type oceanGCPClusterConverter struct{}

func (x *oceanGCPClusterConverter) FromObject(cluster *oceanv1.Cluster) (*OceanCluster, error) {
	return nil, ErrNotImplemented
}

func (x *oceanGCPClusterConverter) ToObject(cluster *OceanCluster) (*oceanv1.Cluster, error) {
	return nil, ErrNotImplemented
}

type oceanGCPLaunchSpecConverter struct{}

func (x *oceanGCPLaunchSpecConverter) FromObject(spec *oceanv1.LaunchSpec) (*OceanLaunchSpec, error) {
	return nil, ErrNotImplemented
}

func (x *oceanGCPLaunchSpecConverter) ToObject(spec *OceanLaunchSpec) (*oceanv1.LaunchSpec, error) {
	return nil, ErrNotImplemented
}
