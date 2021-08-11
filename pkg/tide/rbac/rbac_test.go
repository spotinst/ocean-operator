// Copyright 2021 NetApp, Inc. All Rights Reserved.

package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRBACManifests(t *testing.T) {
	var expectedServiceAccount = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: tide
  namespace: my-namespace
`

	var expectedRoleBinding = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tide-helmadmin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tide
    namespace: my-namespace
`

	t.Run("whenSuccessful", func(tt *testing.T) {
		res, err := GetRBACManifests("my-namespace")
		assert.NoError(tt, err)

		assert.Equal(tt, expectedServiceAccount, res.ServiceAccount)
		assert.Equal(tt, expectedRoleBinding, res.ClusterRoleBinding)
	})
}
