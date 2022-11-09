import type { QueryEditorProps } from '@grafana/data';
import type { NumaflowDataSource } from 'datasource';
import type { NumaflowDataSourceOptions, NumaflowDataQuery } from '../../types';

export type EditorProps = QueryEditorProps<NumaflowDataSource, NumaflowDataQuery, NumaflowDataSourceOptions>;

export type ChangeOptions<T> = {
  propertyName: keyof T;
  runQuery: boolean;
};
