module github.com/spotinst/ocean-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/go-cmp v0.5.5
	github.com/hashicorp/go-version v1.2.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spotinst/spotinst-sdk-go v1.97.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.18.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.22.0
	k8s.io/cli-runtime v0.21.3
	k8s.io/client-go v0.22.0
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.9.5
)
