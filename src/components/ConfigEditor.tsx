import React, { ChangeEvent, PureComponent } from 'react';
import { FieldSet, InlineField, InlineSwitch, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { NumaflowDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<NumaflowDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onNamespacedChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      namespaced: event.target.checked,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onNamespaceChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      namespace: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  render() {
    const { options } = this.props;
    const { jsonData } = options;

    return (
      <>
        <FieldSet label="Kubernetes">
          <InlineField label="Namespaced" tooltip="Whether to run in namespaced scope.">
            <InlineSwitch
              onChange={this.onNamespacedChange}
              placeholder="namespaced"
              value={jsonData?.namespaced ?? false}
            />
          </InlineField>
          <InlineField label="Namespace" tooltip='The namespace to query when "namespaced" is enabled.'>
            <Input onChange={this.onNamespaceChange} placeholder="namespace" value={jsonData?.namespace ?? ''} />
          </InlineField>
        </FieldSet>
      </>
    );
  }
}
