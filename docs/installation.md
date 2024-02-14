# Operator Installation

- [Install on Minikube](#install-on-minikube)
- [Install on Amazon Elastic Kubernetes Service (EKS)](#install-on-amazon-elastic-kubernetes-service-eks)
- [Install on Google Kubernetes Engine (GKE)](#install-on-google-kubernetes-engine-gke)
- [Install on Azure Kubernetes Service (AKS)](#install-on-azure-kubernetes-service-aks)
- [Operator Configuration](#operator-configuration)

## Install on minikube

1. Create minikube cluster.
```bash
minikube start
```
2. Add the helm repo.
```bash
helm repo add metabase-operator-charts https://unagex.github.io/metabase-operator
helm repo update
```
3. Install the operator in the namespace `metabase-operator`. See [operator configuration](#operator-configuration) for more customization.
```bash
helm install metabase-operator metabase-operator-charts/metabase-operator -n metabase-operator --create-namespace
```
Congratulations ! The operator is now installed. To test it, you can deploy a basic metabase (optional):

4. Deploy a metabase in the namespace `default`. See [metabase configuration](./configuration.md) for more customization.
```bash
kubectl apply -f https://raw.githubusercontent.com/unagex/metabase-operator/main/config/samples/v1_metabase.yaml
```
5. Access metabase web UI.
```bash
 minikube service metabase-sample-http
```

## Install on Amazon Elastic Kubernetes Service (EKS)

You should have an EKS cluster already running. See the [official documentation](https://docs.aws.amazon.com/eks/latest/userguide/create-cluster.html) if that's not the case.
1. Install the Amazon EBS CSI driver add-on on your EKS cluster. See the [official documentation](https://docs.aws.amazon.com/eks/latest/userguide/managing-ebs-csi.html) to add it.
2. Grant permissions for your EKS cluster to interact with Amazon EBS volumes, you need to update the IAM roles associated with your EKS nodes. Here is the necessary policy to attach to your cluster role.
```json
{
 "Version": "2012-10-17",
 "Statement": [{
  "Sid": "VisualEditor0",
  "Effect": "Allow",
  "Action": [
   "ec2:CreateVolume",
   "ec2:DeleteVolume",
   "ec2:AttachVolume",
   "ec2:DetachVolume",
   "ec2:DescribeVolumes",
   "ec2:CreateTags",
   "ec2:DeleteTags",
   "ec2:DescribeTags"
  ],
  "Resource": "*"
 }]
}
```
3. Add the helm repo.
```bash
helm repo add metabase-operator-charts https://unagex.github.io/metabase-operator
helm repo update
```
4. Install the operator in the namespace `metabase-operator`. See [operator configuration](#operator-configuration) for more customization.
```bash
helm install metabase-operator metabase-operator-charts/metabase-operator -n metabase-operator --create-namespace
```
Congratulations ! The operator is now installed. To test it, you can deploy a basic metabase (optional):

5. Deploy a metabase in the namespace `default`. See [metabase configuration](./configuration.md) for more customization.
```bash
kubectl apply -f https://raw.githubusercontent.com/unagex/metabase-operator/main/config/samples/v1_metabase.yaml
```
6. Access metabase web UI on port 3000.
```bash
kubectl port-forward services/metabase-sample-http 3000:3000
```

## Install on Google Kubernetes Engine (GKE)

You should have a GKE cluster already running. See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-zonal-cluster) if that's not the case. The following has been tested on a GKE autopilot mode.
1. Add the helm repo.
```bash
helm repo add metabase-operator-charts https://unagex.github.io/metabase-operator
helm repo update
```
2. Install the operator in the namespace `metabase-operator`. See [operator configuration](#operator-configuration) for more customization.
```bash
helm install metabase-operator metabase-operator-charts/metabase-operator -n metabase-operator --create-namespace
```
Congratulations ! The operator is now installed. To test it, you can deploy a basic metabase (optional):

3. Deploy a metabase in the namespace `default`. See [metabase configuration](./configuration.md) for more customization.
```bash
kubectl apply -f https://raw.githubusercontent.com/unagex/metabase-operator/main/config/samples/v1_metabase.yaml
```
4. Access metabase web UI on port 3000.
```bash
kubectl port-forward services/metabase-sample-http 3000:3000
```

## Install on Azure Kubernetes Service (AKS)
You should have a AKS cluster already running. See the [official documentation](https://learn.microsoft.com/en-us/azure/aks/learn/quick-kubernetes-deploy-portal?tabs=azure-cli) if that's not the case.
1. Add the helm repo.
```bash
helm repo add metabase-operator-charts https://unagex.github.io/metabase-operator
helm repo update
```
2. Install the operator in the namespace `metabase-operator`. See [operator configuration](#operator-configuration) for more customization.
```bash
helm install metabase-operator metabase-operator-charts/metabase-operator -n metabase-operator --create-namespace
```
Congratulations ! The operator is now installed. To test it, you can deploy a basic metabase (optional):

3. Deploy a metabase in the namespace `default`. See [metabase configuration](./configuration.md) for more customization.
```bash
kubectl apply -f https://raw.githubusercontent.com/unagex/metabase-operator/main/config/samples/v1_metabase.yaml
```
4. Access metabase web UI on port 3000.
```bash
kubectl port-forward services/metabase-sample-http 3000:3000
```
# Operator configuration

The operator Helm chart is deployed by default with [this values.yaml](/charts/operator/values.yaml). The following values can be overriden:

| Name | Type | Default value
| --- | --- | --- |
| operator.image.repository | string | "ghcr.io/unagex/metabase-operator/controller" |
| operator.image.tag | string | Default to latest version at time of installation. |
| operator.image.pullPolicy | string | "IfNotPresent" |
| resources.limits.cpu | string | nil |
| resources.limits.memory | string | nil |
| resources.requests.cpu | string | nil |
| resources.requests.memory | string | nil |
| labels | map[string]string | nil |
