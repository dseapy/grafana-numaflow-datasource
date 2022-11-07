import type { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NumaflowQuery extends DataQuery {
  rawQuery: string;
}

/**
 * These are options configured for each DataSource instance
 */
export interface NumaflowDataSourceOptions extends DataSourceJsonData {}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */

export type QueryTypesResponse = {
  queryTypes: string[];
};
