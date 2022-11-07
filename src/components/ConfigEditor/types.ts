import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import type { NumaflowDataSourceOptions, NumaflowSecureJsonData } from '../../types';

export type EditorProps = DataSourcePluginOptionsEditorProps<NumaflowDataSourceOptions, NumaflowSecureJsonData>;
