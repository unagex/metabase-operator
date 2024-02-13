# Unagex Kubernetes Operator for Metabase

![Go version](https://img.shields.io/github/go-mod/go-version/unagex/metabase-operator)
![Kubernetes Version](https://img.shields.io/badge/Kubernetes-1.18%2B-green.svg)
![Release](https://img.shields.io/github/v/release/unagex/metabase-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/unagex/metabase-operator)](https://goreportcard.com/report/github.com/unagex/metabase-operator)

## Features


- Create Metabase instances defined as custom resources
- Customize resources needed (cpu, ram) based on your need
- Update Metabase version and config (soon)
- Production-ready with dedicated database

## Quickstart
1. Deploy the operator with helm
```
helm repo add metabase-operator-charts https://unagex.github.io/metabase-operator
helm repo update
helm install metabase-operator metabase-operator-charts/metabase-operator
```
2. Deploy a basic Metabase
```
kubectl apply -f https://raw.githubusercontent.com/unagex/metabase-operator/main/config/samples/v1_metabase.yaml
```
⬇ See documentation below for more ⬇

## Documentation

* [Operator Installation](./docs/installation.md)
* [Metabase Configuration](./docs/configuration.md)
* [Operator Overview](./docs/overview.md)
<br></br>
* [Contribution](./docs/contribution.md)
* [Contact](mailto:mathieu.cesbron@protonmail.com)