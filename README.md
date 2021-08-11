# Ocean Operator for Kubernetes

**Ocean Operator for Kubernetes** is an [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that makes use of [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and can be used to create and manage Ocean components.

## Table of Contents

- [Installation](#installation)
- [Documentation](#documentation)
- [Getting Help](#getting-help)
- [Community](#community)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Install with Spot CLI

#### Prerequisites

- [Install the Spot CLI](https://github.com/spotinst/spotctl#installation).

#### Steps

1. Install `ocean-operator`:

```sh
spotctl ocean install \
  --set spotinst.token=REDACTED \
  --set spotinst.account=REDACTED \
  --set spotinst.clusterIdentifier=REDACTED \
  --namespace spot-system
  # [...]
```

#### Verify

Ensure all Kubernetes Pods in `spot-system` namespace are deployed and have a `STATUS` of `Running`:

```sh
kubectl get pods -n spot-system
```

### Install with Helm

#### Prerequisites

- [Install a Helm client](https://helm.sh/docs/intro/install/) with a version 3 or later.

#### Steps

1. Add the Spot Helm repository:

```sh
helm repo add spot https://charts.spot.io
```

2. Update your local Helm chart repository cache:

```sh
helm repo update
```

3. Install `ocean-operator`:

```sh
helm install my-release spot/ocean-operator \
  --set spotinst.token=REDACTED \
  --set spotinst.account=REDACTED \
  --set spotinst.clusterIdentifier=REDACTED \
  --namespace spot-system \
  --create-namespace \
  # [...]
```

#### Verify

Ensure all Kubernetes Pods in `spot-system` namespace are deployed and have a `STATUS` of `Running`:

```sh
kubectl get pods -n spot-system
```

## Documentation

If you're new to [Spot](https://spot.io/) and want to get started, please checkout our [Getting Started](https://docs.spot.io/connect-your-cloud-provider/) guide, available on the [Spot Documentation](https://docs.spot.io/) website.

## Getting Help

We use GitHub issues for tracking bugs and feature requests. Please use these community resources for getting help:

- Ask a question on [Stack Overflow](https://stackoverflow.com/) and tag it with [ocean-operator](https://stackoverflow.com/questions/tagged/ocean-operator/).
- Join our [Spot](https://spot.io/) community on [Slack](http://slack.spot.io/).
- Open an [issue](https://github.com/spotinst/ocean-operator/issues/new/choose/).

## Community

- [Slack](http://slack.spot.io/)
- [Twitter](https://twitter.com/spot_hq/)

## Contributing

Please see the [contribution guidelines](.github/CONTRIBUTING.md).

## License

Code is licensed under the [Apache License 2.0](LICENSE).
