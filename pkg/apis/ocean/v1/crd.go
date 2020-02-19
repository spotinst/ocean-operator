package v1

import (
	"fmt"

	"gopkg.in/yaml.v3"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ClusterKind is the kind of the Cluster CRD.
	ClusterKind = "Cluster"

	// ClusterSingularName is the singular name of the Cluster CRD.
	ClusterSingularName = "cluster"

	// ClusterPluralName is the plural name of the Cluster CRD.
	ClusterPluralName = "clusters"

	// LaunchSpecPluralName is the kind of the LaunchSpec CRD.
	LaunchSpecKind = "LaunchSpec"

	// LaunchSpecSingularName is the singular name of the LaunchSpec CRD.
	LaunchSpecSingularName = "launchspec"

	// LaunchSpecPluralName is the plural name of the LaunchSpec CRD.
	LaunchSpecPluralName = "launchspecs"
)

// CustomResourceDefinition returns a CustomResourceDefinition object for Cluster.
func (in Cluster) CustomResourceDefinition() *apiextv1beta1.CustomResourceDefinition {
	return MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ClusterPluralName))
}

// CustomResourceDefinition returns a CustomResourceDefinition object for LaunchSpec.
func (in LaunchSpec) CustomResourceDefinition() *apiextv1beta1.CustomResourceDefinition {
	return MustCustomResourceDefinition(SchemeGroupVersion.WithResource(LaunchSpecPluralName))
}

// CustomResourceDefinition returns a CustomResourceDefinition from the asset
// packaged by go-bindata, or an error in case of failure.
func CustomResourceDefinition(gvr schema.GroupVersionResource) (*apiextv1beta1.CustomResourceDefinition, error) {
	data, err := Asset(fmt.Sprintf("%s_%s_crd.yaml", gvr.Group, gvr.Resource))
	if err != nil {
		return nil, err
	}

	var out apiextv1beta1.CustomResourceDefinition
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	out.ObjectMeta.Name = fmt.Sprintf("%s.%s", gvr.Resource, gvr.Group)
	out.Spec.Group = gvr.Group

	return &out, nil
}

// MustCustomResourceDefinition returns a CustomResourceDefinition from the asset
// packaged by go-bindata and panics in case of failure.
func MustCustomResourceDefinition(gvr schema.GroupVersionResource) *apiextv1beta1.CustomResourceDefinition {
	out, err := CustomResourceDefinition(gvr)
	if err != nil {
		panic(err)
	}
	return out
}
