# [[project-name]]

Chart version 0.1.0, app version 0.1.0

A Helm chart to deploy the sample golang based web service project.

**Homepage:** <https://github.com/[[repo-owner]]/[[image-name]]/>

## Requirements

Kubernetes: `>=1.25.0`

## Values

Use the values below to configure the chart's values.
| Key | Type | Default | Description |
|-----|------|---------|-------------|
| autoscaling.enabled | bool | `true` | Turn on Pod replicas number autoscaling instead of setting a constant value. your cluster must support [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/). |
| autoscaling.maxReplicas | int | `10` | Max number of Pods autoscaler can deploy. |
| autoscaling.minReplicas | int | `3` | Min number of Pods autoscaler can deploy. |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Pod scale up critieria based on CPU usage. |
| autoscaling.targetMemoryUtilizationPercentage | int | `80` | Pod scale up critieria based on Memory usage. |
| database | object | `{"connectionSecretName":"","name":"[[project-name]]"}` | Secret that contains database connection details. It has to be present in the same namespace as the Chart is deployed to. Has to have `endpoint`, `username` and `password` keys. May contain `port` key. |
| database.connectionSecretName | string | `""` | Secret that contains database connection details. It has to be present in the same namespace as the Chart is deployed to. Has to have `endpoint`, `username` and `password` keys. May contain `port` key. |
| database.name | string | `"[[project-name]]"` | Name of a logical database to use |
| fullnameOverride | string | `""` | Override the default name generated for this specific chart Release. |
| ginMode | string | `"debug"` | Configure run mode of the gin web framework; can be 'debug' or 'production' |
| image.pullPolicy | string | `"IfNotPresent"` | Configure image pull policy. |
| image.registry | string | `"[[registry-domain]]"` | Set the domain of your container images registry. |
| image.repository | string | `"[[registry-name]]/[[image-name]]"` | Set the name of the repository within the registry. |
| image.tag | string | `""` | Image tag to use, defaults to .Chart.AppVersion |
| imagePullSecrets | list | `[]` | Configure login secrets for the container images registry. |
| inMemoryStore | bool | `false` | When set to "true", the app doesn't use any database at all and stores all the data in memory only. |
| ingress.annotations | object | `{}` | Optional annotations for the Ingress definition. If your cluster has "CertManager" operator running, you can use "cert-manager.io/cluster-issuer" annotation to [automatically generate a certificate for it](https://cert-manager.io/docs/usage/). |
| ingress.className | string | `"nginx"` | Ingress controller implementations use this field to know whether they should be serving this Ingress resource, by a transitive connection. |
| ingress.enabled | bool | `true` | Should the Service be accessible through an Ingress. This needs an Ingress controller to be configured already on your cluster. |
| ingress.host | string | `"chart-example.local"` | HTTP host that you want to use for your service. |
| ingress.tls | list | `[]` | Optional TLS certificate configuration. You can use it with "CertManager" or provide your own certificate. |
| monitoring.serviceMonitor | object | `{"enabled":true,"extraLabels":{}}` | If your cluster supports prometheus-operator configuration of metrics data, enable this to have metrics from your application automatically ingested by prometheus. |
| monitoring.serviceMonitor.extraLabels | object | `{}` | Optional extra labels to put on the serviceMonitor |
| nameOverride | string | `""` | Override the default name generated for the chart objects. |
| nodeSelector | object | `{}` | Optional node delector to limit the nodes where pods of the chart can be deployed. |
| pdb | object | `{"enabled":true}` | Should the chart deploy a [PodDisruptionBudget](https://kubernetes.io/docs/tasks/run-application/configure-pdb/) to limit disruptions based on administrative tasks. |
| podAnnotations | object | `{}` | Set additional annotations for the pods created. |
| podListenPort | int | `8080` | Configure the TCP port on which your pods will listen for connections. |
| redisConnectionSecretName | string | `""` | Secret that contains redis connection details. It has to be present in the same namespace as the Chart is deployed to. Has to have `host`, `username` and `password` keys. May contain `port` key. |
| replicaCount | int | `3` | Number of Pod replicas to deploy. Used only if 'autoscaling.enabled' is 'false'. |
| resources | object | `{"limits":{"cpu":"500m","memory":"512Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}` | Configure [Pod resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/). |
| service.annotations | object | `{}` | Optional annotations for the Service definition. If your cluster has "ExternalDNS" operator running, you can use "external-dns.alpha.kubernetes.io/hostname" annotation to [automatically register DNS name for your service](https://github.com/kubernetes-sigs/external-dns). |
| service.port | int | `80` | TCP port that the service will be exposed on. |
| service.type | string | `"ClusterIP"` | The type of [Service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types). |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account. |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created. |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |

## Maintainers

| Name           | Email                 | Url |
| -------------- | --------------------- | --- |
| [[repo-owner]] | <noemail@nothing.com> |     |
