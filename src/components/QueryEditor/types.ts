import type { QueryEditorProps } from '@grafana/data';
import type { NumaflowDataSource } from 'datasource';
import type { NumaflowDataSourceOptions, NumaflowQuery } from '../../types';

export type EditorProps = QueryEditorProps<NumaflowDataSource, NumaflowQuery, NumaflowDataSourceOptions>;

export type ChangeOptions<T> = {
  propertyName: keyof T;
  runQuery: boolean;
};
