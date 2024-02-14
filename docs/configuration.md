# Metabase Configuration

After deploying the operator, you can create a `Metabase` resource to create a Metabase instance ([example Metabase resource](../config/samples/v1_metabase.yaml)).

## Metabase Spec

The following `spec values` can be updated.

| Name | Type | Default value | Description |
| --- | --- | --- | --- |
| metabase.image | string | "metabase/metabase:latest" | The Metabase image. |
| metabase.imagePullPolicy | string | "IfNotPresent" | ImagePullPolicy of Metabase image. |
| metabase.resources | corev1.ResourceRequirements | {requests: {cpu: "1", memory: "2Gi"}, limits: {cpu: "1", memory: "2Gi"}} | Resources of Metabase container. |
| db.image | string | "postgres:latest" | The PostgreSQL image. |
| db.imagePullPolicy | "IfNotPresent" | ImagePullPolicy of DB image. |
| db.replicas | 1 | Number of replicas of PostgreSQL image. |
| db.resources | corev1.ResourceRequirements | {requests: {cpu: "100m", memory: "256Mi"}, limits: {cpu: "1", memory: "2Gi"}} | Resources of PostgreSQL container. |
| db.volume.storageClassName | string | Name of the default storageClass of the cluster.  | StorageClassName of the PostgreSQL. |
| db.volume.size | string | "10Gi" | Size of the database, e.g., 50Mi, 2Gi. |

## Metabase Status

The `status` field of the `Metabase` resource is updated by the Metabase operator. It gives real-time status.

| Name | Type | Description |
| --- | --- | --- |
| ready | bool | The Metabase instance can be used when the value is `true`.|
| hosts.HTTP | string | Host to access Metabase instance.|