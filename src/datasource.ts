import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import type { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import type { NumaflowQuery, NumaflowDataSourceOptions, QueryTypesResponse } from './types';

export class NumaflowDataSource extends DataSourceWithBackend<NumaflowQuery, NumaflowDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NumaflowDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: NumaflowQuery, scopedVars: ScopedVars): Record<string, any> {
    return {
      ...query,
      rawQuery: getTemplateSrv().replace(query.rawQuery, scopedVars),
    };
  }

  getAvailableQueryTypes(): Promise<QueryTypesResponse> {
    return this.getResource('/query-types');
  }
}
