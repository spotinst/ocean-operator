#!/bin/bash
#
# Copyright 2020 Spotinst LTD.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

B64_COMMAND="base64 --wrap=0"
B64_SPOTINST_TOKEN="$(echo -n ${SPOTINST_TOKEN} | ${B64_COMMAND})"
B64_SPOTINST_ACCOUNT="$(echo -n ${SPOTINST_ACCOUNT} | ${B64_COMMAND})"

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: ocean-operator
type: Opaque
data:
  token: ${B64_SPOTINST_TOKEN}
  account: ${B64_SPOTINST_ACCOUNT}
EOF

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: ocean-operator
data:
  config.yaml: |-
    # Bootstrap configuration.
    bootstrap:

      # List of Custom Resource Definitions to install.
      customResourceDefinitions:

      # Cluster.
      - name: clusters.ocean.spot.io
        installPolicy: Always

      # Launch Spec.
      - name: launchspecs.ocean.spot.io
        installPolicy: Always
EOF
