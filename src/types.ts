import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NumaflowDataQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const defaultQuery: Partial<NumaflowDataQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface NumaflowDataSourceOptions extends DataSourceJsonData {
  namespaced?: boolean;
  namespace?: string;
}
