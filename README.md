# Grafana Backend Datasource Plugin

The following assumes you are using the [Grafana Helm Chart](https://github.com/grafana/helm-charts/tree/main/charts/grafana)
either directly or indirectly (i.e. [Kube Prometheus Stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)).

If installing Grafana through some other means, the process should be similar.

## Build
```shell
docker build -t <my-container-registry>/grafana-numaflow-datasource:latest .
docker push <my-container-registry>/grafana-numaflow-datasource:latest
```

## Configure Grafana

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

## Configure Datasource

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

## Configure Dashboards

TODO