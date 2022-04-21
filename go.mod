module github.com/spotinst/ocean-operator

go 1.16

require (
	github.com/go-logr/logr v1.2.0
	github.com/google/go-cmp v0.5.6
	github.com/hashicorp/go-version v1.3.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spotinst/spotinst-sdk-go v1.105.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.7.1
	k8s.io/api v0.23.6
	k8s.io/apiextensions-apiserver v0.22.4
	k8s.io/apimachinery v0.23.6
	k8s.io/cli-runtime v0.23.6
	k8s.io/client-go v0.23.6
	sigs.k8s.io/controller-runtime v0.10.3
)
