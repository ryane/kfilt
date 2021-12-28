# kfilt

[![Build Status](https://travis-ci.org/ryane/kfilt.svg?branch=master)](https://travis-ci.org/ryane/kfilt)
[![Code Coverage](https://codecov.io/gh/ryane/kfilt/branch/master/graph/badge.svg)](https://codecov.io/gh/ryane/kfilt)
[![Go Report Card](https://goreportcard.com/badge/ryane/kfilt)](https://goreportcard.com/report/ryane/kfilt)
[![LICENSE](https://img.shields.io/github/license/ryane/kfilt.svg)](https://github.com/ryane/kfilt/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/release-pre/ryane/kfilt.svg)](https://github.com/ryane/kfilt/releases)

kfilt can filter Kubernetes resources.

## What is kfilt?

kfilt is a tool that lets you filter specific resources from a stream of Kubernetes YAML manifests. It can read manifests from a file, URL, or from stdin.

kfilt was primarily created to assist developers who are creating [Helm charts](https://helm.sh/docs/developing_charts/) or [Kustomize](https://github.com/kubernetes-sigs/kustomize) bases. Often, when making changes, it is helpful to narrow down focus to a specific resource or set of resources in the output. Without kfilt, you might redirect output to a file for inspection in your text editor or to write complicated grep commands. kfilt makes it easy to filter the output to see just the resources you are currently interested in. Or, to exclude specific resources.

You can also use kfilt to selectively apply (or delete) resources with kubectl.

It is easiest to understand with some examples.

### Examples

#### Working with Files or URLs

Only output a ServiceAccount named "test":

```
kfilt -f ./pkg/decoder/test.yaml -i kind=ServiceAccount,name=test
```

Output all Service Accounts:

```
kfilt -f http://bit.ly/2xSiCJL -i kind=ServiceAccount
```

#### Working with Helm charts

kfilt can be used as a more flexible alternative to the `-x` option of [`helm template`](https://helm.sh/docs/helm/#helm-template). In this example, we will only output rendered Service resources from a Helm chart:

```
helm template chart | kfilt -i kind=service
```

We also have the ability to exclude resources. Here we will exclude all Secrets before applying a Chart to a cluster.

```
helm template chart | kfilt -x kind=secret | kubectl apply -f -
```

#### Working with Kustomize

Only output the ConfigMaps in a Kustomize base.

```
kustomize build github.com/kubernetes-sigs/kustomize//examples/helloWorld | kfilt -i kind=ConfigMap
```

Output only the resources named "the-deployment".

```
kustomize build github.com/kubernetes-sigs/kustomize//examples/helloWorld | kfilt -i name=the-deployment
```

#### Working with kubectl

Find all resources named "nginx-ingress-controller" regardless of kind.

```
kubectl get all -A -oyaml | kfilt -i name=nginx-ingress-controller
```

## Installation

kfilt is available on Linux, Mac, and Windows <sup>1</sup> and binaries are available on the [releases](https://github.com/ryane/kfilt/releases) page.

### Docker

You can also run kfilt as a Docker container. Make sure you include `-i` in your `docker run` command.

```
kustomize build base | docker run --rm -i ryane/kfilt -k ConfigMap
```

### Running as a Kustomize Plugin (experimental)

See [plugin/kustomize](./plugin/kustomize) for an experimental Kustomize plugin.

## Usage

### Including Resources

You can use `--include` or `-i` to control which resources to include in the kfilt output. This argument takes a list of simple key value pairs that make up your query. The following keys are currently supported:

| Key           | Field              | Example                   |
|---------------|--------------------|---------------------------|
| kind, k       | kind               | ServiceAccount            |
| name, n       | metadata.name      | my-app                    |
| group, g      | apiVersion         | rbac.authorization.k8s.io |
| version, v    | apiVersion         | v1                        |
| namespace, ns | metadata.namespace | kube-system               |
| labels, l     | metadata.labels    | app=my-app                |

Note that it is possible to use wildcards (`*`, and `?`) when filtering by name.

#### Examples

##### Filter by kind

```
kfilt -f ./pkg/decoder/test.yaml -i kind=configmap
```

##### Filter by group

```
kfilt -f ./pkg/decoder/test.yaml -i g=config.istio.io
```

##### Filter by name and kind

You can combine keys in a single `--include` by separating them with a comma. In this example, we are filtering to match ServiceAccount resources named "test":

```
kfilt -f ./pkg/decoder/test.yaml -i k=ServiceAccount,n=test
```

##### Filter by name using wildcards

```
kfilt -f ./pkg/decoder/test.yaml -n "test*"
```

You can use `*` and `?` wildcard characters.

##### Filter by multiple kinds

You can use multiple `--include` flags. kfilt will output resources that match any one of the includes. For example, to output ServiceAccounts and ConfigMaps you could use:

```
kfilt -f ./pkg/decoder/test.yaml -i k=serviceaccount -i k=configmap
```

##### Filter with Label Selectors

```
kfilt -f ./pkg/decoder/test.yaml -i labels=app=test
```

### Excluding Resources

The `--exclude` or `-x` flag will allow you to exclude resources. This supports the same key value pairs as the `--include` flag.

#### Examples

##### Exclude by kind

```
kfilt -f ./pkg/decoder/test.yaml -x kind=configmap
```

##### Exclude by name

```
kfilt -f ./pkg/decoder/test.yaml -x name=test
```

##### Exclude multiple kinds

```
kfilt -f ./pkg/decoder/test.yaml -x kind=configmap -x k=serviceaccount
```

##### Exclude with Label Selectors

```
kfilt -f ./pkg/decoder/test.yaml -x labels=app=test
```

### Shortcuts

Because "kind", "name", and "labels" are the most commonly used fields to filter by, kfilt has special flags allowing you to save some typing.

You can include by "kind" by using the `--kind` (or `-k`) flag with just the name of the kind you want to filter by. You can use `--exclude-kind` (or `-K`) for exclusions.

The corresponding flags for "name" queries are `--name` (`-n`) and `--exclude-name` (`-N`).

Finally, you can use label selectors with the `--labels` (`-l`) and `--exclude-labels` (`L`) flags.

#### Include ConfigMaps and Service Accounts

```
kfilt -f ./pkg/decoder/test.yaml -k configmap -k serviceaccount
```

#### Exclude resources named "test"

```
kfilt -f ./pkg/decoder/test.yaml -N test
```

#### Include resources labeled with "app=test"

```
kfilt -f ./pkg/decoder/test.yaml -l app=test
```

---

<sup>1</sup> *note*: kfilt has not been tested extensively on Windows. Please file an issue if you run into any problems.
