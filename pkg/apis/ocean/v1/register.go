// NOTE: Boilerplate only.  Ignore this file.

// Package v1 contains API Schema definitions for the ocean v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=ocean.spot.io
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "ocean.spot.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)
