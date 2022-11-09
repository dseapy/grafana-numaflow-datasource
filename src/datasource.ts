import {DataSourceWithBackend, getTemplateSrv} from '@grafana/runtime';
import type {DataSourceInstanceSettings, ScopedVars} from '@grafana/data';
import {MetricNamesResponse, NumaflowDataQuery, NumaflowDataSourceOptions, QueryTypesResponse} from "./types";

export class NumaflowDataSource extends DataSourceWithBackend<NumaflowDataQuery, NumaflowDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NumaflowDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: NumaflowDataQuery, scopedVars: ScopedVars): Record<string, any> {
    return {
      ...query,
      rawQuery: getTemplateSrv().replace(query.rawQuery, scopedVars),
    };
  }

  getAvailableQueryTypes(): Promise<QueryTypesResponse> {
    return this.getResource('/query-types');
  }

  fetchMetricNames(query: string): Promise<MetricNamesResponse> {
    return this.postResource('/metric-names', { rawQuery: query } )
  }

async metricFindQuery(query: string, options?: any) {
    const response = await this.fetchMetricNames(query);
    return response.metricNames.map(name => ({text: name}));
  }
}
