import { DataSourcePlugin } from '@grafana/data';
import { NumaflowDataSource } from './datasource';
import type { NumaflowDataQuery, NumaflowDataSourceOptions } from './types';
import { ConfigEditor, QueryEditor } from './components';
import {VariableQueryEditor} from "./components/QueryEditor/VariableQueryEditor";

export const plugin = new DataSourcePlugin<NumaflowDataSource, NumaflowDataQuery, NumaflowDataSourceOptions>(NumaflowDataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setVariableQueryEditor(VariableQueryEditor);
