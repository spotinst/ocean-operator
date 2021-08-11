// Copyright 2021 NetApp, Inc. All Rights Reserved.

package helm

import (
	"testing"

	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	"github.com/spotinst/ocean-operator/pkg/installer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestYamlConversion(t *testing.T) {
	values := `
serviceAccount:
  create: true
`
	var vals map[string]interface{}
	err := yaml.Unmarshal([]byte(values), &vals)
	assert.NoError(t, err)

	for k, v := range vals {
		t.Log("values", k, v)
	}

	e := vals["serviceAccount"]
	s, ok := e.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, true, s["create"])
}

func getVersionedObjects(componentVersion, releasedVersion string) (*oceanv1alpha1.OceanComponent, *installer.Release) {
	return &oceanv1alpha1.OceanComponent{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ocean-foo",
			},
			Spec: oceanv1alpha1.OceanComponentSpec{
				Name:    "foo",
				Version: componentVersion,
				Values:  "",
			},
		},
		&installer.Release{
			Name:    "foo",
			Version: releasedVersion,
		}
}

func getValuesObjects(componentValues string, releasedValues map[string]interface{}) (*oceanv1alpha1.OceanComponent, *installer.Release) {
	return &oceanv1alpha1.OceanComponent{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ocean-foo",
			},
			Spec: oceanv1alpha1.OceanComponentSpec{
				Name:    "foo",
				Version: "v1.2",
				Values:  componentValues,
			},
		},
		&installer.Release{
			Name:    "foo",
			Version: "v1.2",
			Values:  releasedValues,
		}

}

func TestIsUpgrade(t *testing.T) {
	logger := zap.New(zap.UseDevMode(true)).WithValues("test", t.Name())
	i := &Installer{
		ClientGetter: nil,
		Log:          logger,
	} // fix getClient for more complex tests
	var u bool

	u = i.IsUpgrade(getVersionedObjects("v1.1.0", "v0.9.8"))
	assert.True(t, u)

	u = i.IsUpgrade(getVersionedObjects("v1.1.0", "v1.1.0"))
	assert.False(t, u)

	u = i.IsUpgrade(getValuesObjects("metricsEnabled: true", map[string]interface{}{}))
	assert.True(t, u)

	u = i.IsUpgrade(getValuesObjects("", map[string]interface{}{}))
	assert.False(t, u)

	u = i.IsUpgrade(getValuesObjects(":unparseable yaml is an upgrade lol:", map[string]interface{}{}))
	assert.True(t, u)

	v1 := `
serviceAccount:
  create: true
`
	v2 := map[string]interface{}{
		"serviceAccount": map[string]interface{}{
			"create": true,
		},
	}
	u = i.IsUpgrade(getValuesObjects(v1, v2))
	assert.False(t, u)

}
