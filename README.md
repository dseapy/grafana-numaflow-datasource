# Grafana Numaflow Datasource Plugin

A Numaflow [backend data source plugin](https://grafana.com/tutorials/build-a-data-source-backend-plugin/) for Grafana

Supports viewing Numaflow metadata & metrics that come from the following sources:
* Kubernetes Cluster Resources (i.e. `Pipeline`, `Vertex`, `InterStepBufferService`)
* Numaflow Daemon Services
* Kubernetes Metrics API

Supports the following Grafana panels:
* [Table](https://grafana.com/docs/grafana/v9.0/visualizations/table/)
* [Node Graph](https://grafana.com/docs/grafana/v9.0/visualizations/node-graph/) (Only supported when querying for a single `Pipeline`)

**Still a proof-of-concept** - code rushed, no unit tests, etc.

Read more here: https://medium.com/@dseapy/monitoring-stream-processing-in-a-kubernetes-native-environment-9f8f68e82346

## Build
```shell
docker build -t <my-container-registry>/grafana-numaflow-datasource:latest .
docker push <my-container-registry>/grafana-numaflow-datasource:latest
```

## Install

The following assumes you are using the [Grafana Helm Chart](https://github.com/grafana/helm-charts/tree/main/charts/grafana)
either directly or indirectly (i.e. [Kube Prometheus Stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)).

If installing Grafana through some other means, the process should be similar.

### Grafana

Add the following helm overrides to your Grafana installation.

```yaml
grafana.ini:
  # Consider signing and remove this
  plugins:
    allow_loading_unsigned_plugins: "numaflow-datasource"
  # (Optional) Can increase log level to debug if necessary
  log:
    level: debug
extraInitContainers:
- name: init-plugins
  image: <my-container-registry>/grafana-numaflow-datasource:latest
  command: [ "/bin/sh", "-c", "mkdir /plugins/numaflow && cp -r /dist/* /plugins/numaflow/"]
  resources:
    limits:
      memory: 100Mi
    requests:
      cpu: 50m
  volumeMounts:
  - name: plugins
    mountPath: /plugins
extraContainerVolumes:
- name: plugins
  emptyDir: {}
extraVolumeMounts:
- name: plugins
  mountPath: /var/lib/grafana/plugins
sidecar:
  resources:
    limits:
      memory: 100Mi
    requests:
      cpu: 50m
  dashboards:
     enabled: true
     searchNamespace: ALL
  datasources:
    enabled: true
    searchNamespace: ALL
```

### Datasource

Add the following `ConfigMap` to your Numaflow installation.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-numaflow-datasource
  labels:
    grafana_datasource: "1"
data:
  datasource.yaml: |-
    apiVersion: 1
    datasources:
      - name: Numaflow
        type: numaflow-datasource
        uid: numaflow
        jsonData:
          namespaced: false
```

### Dashboards

TODO

## Queries

The following assumes you are using variables `$namespace`, `$pipeline`, `$vertex`, `$isbsvc` in grafana.
Replace with other values (i.e. `my-namespace`) if not using Grafana variables.

### Metric Names (for variables)
All pipelines in namespace:
```json
{"namespace":"$namespace","pipeline":"*"}
```
All vertices in namespace
```json
{"namespace":"$namespace","pipeline":"","vertex":"*"}
```
All vertices in pipeline
```json
{"namespace":"$namespace","pipeline":"$pipeline","vertex":"*"}
```
All pods in vertex
```json
{"namespace":"$namespace","pipeline":"$pipeline","vertex":"$vertex","pod":"*"}
```
All pods in isbsvc
```json
{"namespace":"$namespace","isbsvc":"$isbsvc","pod":"*"}
```
All isbsvcs in namespace
```json
{"namespace":"$namespace","isbsvc":"*"}
```
All namespaces containing pipelines
```json
{"namespace":"*","pipeline":""}
```
All namespaces containing vertices
```json
{"namespace":"*","pipeline":"","vertex":""}
```
All namespaces containing isbsvcs
```json
{"namespace":"*","isbsvc":""}
```

### Data (for Table and NodeGraph panels)
All pipelines in all namespaces:
```json
{"namespace":"","pipeline":"*"}
```
All vertices in all pipelines:
```json
{"namespace":"","pipeline":"","vertex":"*"}
```
All isbsvcs in all namespaces
```json
{"namespace":"","isbsvc":"*"}
```
All pipelines in namespace
```json
{"namespace":"$namespace","pipeline":"*"}
```
All vertices in namespace
```json
{"namespace":"$namespace","pipeline":"","vertex":"*"}
```
All vertices in pipeline
```json
{"namespace":"$namespace","pipeline":"$pipeline","vertex":"*"}
```
All isbsvcs in namespace
```json
{"namespace":"$namespace","isbsvc":"*"}
```
A single pipeline
```json
{"namespace":"$namespace","pipeline":"$pipeline"}
```
A single vertex
```json
{"namespace":"$namespace","pipeline":"$pipeline","vertex":"$vertex"}
```
A single isbsvc
```json
{"namespace":"$namespace","isbsvc":"$isbsvc"}
```
