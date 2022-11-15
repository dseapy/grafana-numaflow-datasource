import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import type { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { MetricNamesResponse, NumaflowDataQuery, NumaflowDataSourceOptions, QueryTypesResponse } from './types';
import { MultiValueVariable, TextValuePair } from './components/QueryEditor/types';
import _ from 'lodash';

const supportedVariableTypes = ['constant', 'custom', 'query', 'textbox'];

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
    return this.postResource('/metric-names', { rawQuery: query });
  }

  async metricFindQuery(query: string, options?: any) {
    let payload = query;
    payload = getTemplateSrv().replace(payload, { ...this.getVariables });
    const response = await this.fetchMetricNames(payload);
    return response.metricNames.map((name) => ({ text: name }));
  }

  getVariables() {
    const variables: { [id: string]: TextValuePair } = {};
    Object.values(getTemplateSrv().getVariables()).forEach((variable) => {
      if (!supportedVariableTypes.includes(variable.type)) {
        console.warn(`Variable of type "${variable.type}" is not supported`);

        return;
      }

      const supportedVariable = variable as MultiValueVariable;

      let variableValue = supportedVariable.current.value;
      if (variableValue === '$__all' || _.isEqual(variableValue, ['$__all'])) {
        if (supportedVariable.allValue === null || supportedVariable.allValue === '') {
          variableValue = supportedVariable.options.slice(1).map((textValuePair) => textValuePair.value);
        } else {
          variableValue = supportedVariable.allValue;
        }
      }

      variables[supportedVariable.id] = {
        text: supportedVariable.current.text,
        value: variableValue,
      };
    });

    return variables;
  }
}
