import type { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NumaflowDataQuery extends DataQuery {
  rawQuery: string;
}

export interface NumaflowDataSourceOptions extends DataSourceJsonData {
  namespaced?: boolean;
  namespace?: string;
}

export type QueryTypesResponse = {
  queryTypes: string[];
};

export type MetricNamesResponse = {
  metricNames: string[];
};
