import { DataSourcePlugin } from '@grafana/data';
import { NumaflowDataSource } from './datasource';
import type { NumaflowQuery, NumaflowDataSourceOptions } from './types';
import { ConfigEditor, QueryEditor } from './components';

export const plugin = new DataSourcePlugin<NumaflowDataSource, NumaflowQuery, NumaflowDataSourceOptions>(NumaflowDataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
