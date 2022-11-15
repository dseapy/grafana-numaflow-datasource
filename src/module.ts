import { DataSourcePlugin } from '@grafana/data';
import { NumaflowDataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { NumaflowDataQuery, NumaflowDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<NumaflowDataSource, NumaflowDataQuery, NumaflowDataSourceOptions>(NumaflowDataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
