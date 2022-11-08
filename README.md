# Grafana Data Source Plugin for Numaflow

[Grafana data source plugin](https://grafana.com/tutorials/build-a-data-source-plugin/) for [Numaflow](https://github.com/numaproj/numaflow). 

* Use [Prometheus data source](https://grafana.com/docs/grafana/latest/datasources/prometheus/) for the time-series metrics. See Numaflow's [metrics.md](https://github.com/numaproj/numaflow/blob/main/docs/metrics/metrics.md).

* Use a Grafana logging data source (i.e. [Loki](https://grafana.com/docs/grafana/latest/datasources/loki/), [ElasticSearch](https://grafana.com/docs/grafana/latest/datasources/elasticsearch/)) for the container logs.

* Use this data source for metadata in Kubernetes/Numaflow that isn't easily collected from existing grafana data sources.
An example is `edge` & `vertex` metadata for a [node graph panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/node-graph/).

**Disclaimers**:
* Grafana's node graph panel is in beta
* Plugin relies on Numaflow APIs that may change
* Not an official plugin, not affiliated with numaflow, no support guarantees :)

## Resources

- [QUICK_START](docs/quick-start.md)
- [DEVELOPMENT](docs/development.md)
