import React, {ReactElement} from 'react';
import {FieldSet, InlineField, InlineSwitch, Input} from '@grafana/ui';
import type { EditorProps } from './types';
import { useChangeOptions } from './useChangeOptions';
import { testIds } from '../testIds';
import {useChangeSwitch} from "./useChangeSwitch";

export function ConfigEditor(props: EditorProps): ReactElement {
  const { jsonData } = props.options;
  const onNamespacedChange = useChangeSwitch(props, 'namespaced');
  const onNamespaceChange = useChangeOptions(props, 'namespace');


  return (
      <>
        <FieldSet label="Kubernetes">
          <InlineField label="Namespaced" tooltip="Whether to run in namespaced scope.">
            <InlineSwitch
                onChange={onNamespacedChange}
                placeholder="namespaced"
                data-testid={testIds.configEditor.namespaced}
                value={jsonData?.namespaced ?? false}
            />
          </InlineField>
          <InlineField label="Namespace" tooltip='The namespace to query when "namespaced" is "true".'>
            <Input
                onChange={onNamespaceChange}
                placeholder="namespace"
                data-testid={testIds.configEditor.namespace}
                value={jsonData?.namespace ?? ''}
            />
          </InlineField>
        </FieldSet>
      </>
  );
}
