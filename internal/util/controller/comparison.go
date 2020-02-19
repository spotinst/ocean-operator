package controller

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	oceanv1 "github.com/spotinst/ocean-operator/pkg/apis/ocean/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CompareClusterSpecs(a, b oceanv1.ClusterSpec) Comparison {
	return newComparison(Diff(a, b))
}

func CompareClusterStatuses(a, b oceanv1.ClusterStatus) Comparison {
	return newComparison(Diff(a, b))
}

func CompareLaunchSpecSpecs(a, b oceanv1.LaunchSpecSpec) Comparison {
	return newComparison(Diff(a, b))
}

func CompareLaunchSpecStatuses(a, b oceanv1.LaunchSpecStatus) Comparison {
	return newComparison(Diff(a, b))
}

// Diff returns the difference between two objects.
func Diff(a, b interface{}) string {
	return cmp.Diff(a, b, defaultComparisonOptions...)
}

// Equal checks that two objects are equal.
func Equal(a, b interface{}) bool {
	return cmp.Equal(a, b, defaultComparisonOptions...)
}

// Simple struct representing whether two objects match and their differences.
type Comparison struct {
	// A human-readable list of differences.
	Diff string

	// Whether or not the two objects are equal.
	Equal bool
}

// Simple helper function that creates a Comparison object from a difference result.
func newComparison(diff string) Comparison {
	return Comparison{
		Diff:  diff,
		Equal: diff == "",
	}
}

// Default comparison options for specific behavior of Equal and Diff.
var defaultComparisonOptions = []cmp.Option{
	equateEmptySlicesAndMapsToNil,
	equateNilBoolToFalse,
	equateNilStringToEmptyString,
	ignoreTypeMeta,
	ignoreObjectMeta,
	ignoreTagSliceOrder,
	ignoreLabelSliceOrder,
	ignoreTaintSliceOrder,
}

// Equate bool pointers that are nil to false.
var equateNilBoolToFalse = cmp.Transformer("", func(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
})

// Equate string pointers that are nil to the empty string.
var equateNilStringToEmptyString = cmp.Transformer("", func(value *string) string {
	if value == nil {
		return ""
	}
	return *value
})

// Equate slices and maps that are empty to nil.
var equateEmptySlicesAndMapsToNil = cmpopts.EquateEmpty()

// Ignore TypeMeta.
var ignoreTypeMeta = cmpopts.IgnoreTypes(metav1.TypeMeta{})

// Ignore ResourceVersion.
var ignoreObjectMeta = cmpopts.IgnoreFields(metav1.ObjectMeta{}, "ResourceVersion")

// Define a custom sort order for *oceanv1.Tag slices.
var ignoreTagSliceOrder = cmpopts.SortSlices(func(left, right *oceanv1.Tag) bool {
	// Order is arbitrary, but must be deterministic, irreflexive, and transitive.
	// https://godoc.org/github.com/google/go-cmp/cmp/cmpopts#SortSlices
	if left == nil && right == nil {
		return true
	} else if left != nil && right == nil {
		return true
	} else if left == nil && right != nil {
		return false
	}
	if left.Key != right.Key {
		return left.Key < right.Key
	} else {
		return left.Value < right.Value
	}
})

// Define a custom sort order for *oceanv1.Label slices.
var ignoreLabelSliceOrder = cmpopts.SortSlices(func(left, right *oceanv1.Label) bool {
	// Order is arbitrary, but must be deterministic, irreflexive, and transitive.
	// https://godoc.org/github.com/google/go-cmp/cmp/cmpopts#SortSlices
	if left == nil && right == nil {
		return true
	} else if left != nil && right == nil {
		return true
	} else if left == nil && right != nil {
		return false
	}
	if left.Key != right.Key {
		return left.Key < right.Key
	} else {
		return left.Value < right.Value
	}
})

// Define a custom sort order for *oceanv1.Taint slices.
var ignoreTaintSliceOrder = cmpopts.SortSlices(func(left, right *oceanv1.Taint) bool {
	// Order is arbitrary, but must be deterministic, irreflexive, and transitive.
	// https://godoc.org/github.com/google/go-cmp/cmp/cmpopts#SortSlices
	if left == nil && right == nil {
		return true
	} else if left != nil && right == nil {
		return true
	} else if left == nil && right != nil {
		return false
	}
	if left.Key != right.Key {
		return left.Key < right.Key
	} else {
		return left.Value < right.Value
	}
})
