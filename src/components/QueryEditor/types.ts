import type { QueryEditorProps, VariableModel } from '@grafana/data';
import type { NumaflowDataSource } from 'datasource';
import type { NumaflowDataSourceOptions, NumaflowDataQuery } from '../../types';

export type EditorProps = QueryEditorProps<NumaflowDataSource, NumaflowDataQuery, NumaflowDataSourceOptions>;

export type ChangeOptions<T> = {
  propertyName: keyof T;
  runQuery: boolean;
};

export interface TextValuePair {
  text: string;
  value: any;
}

export interface MultiValueVariable extends VariableModel {
  allValue: string | null;
  id: string;
  current: TextValuePair;
  options: TextValuePair[];
}