import { ChangeEvent, useCallback } from 'react';
import type { NumaflowDataSourceOptions } from 'types';
import type { EditorProps } from './types';

type OnChangeType = (event: ChangeEvent<HTMLInputElement>) => void;

export function useChangeSwitch(props: EditorProps, propertyName: keyof NumaflowDataSourceOptions): OnChangeType {
    const { onOptionsChange, options } = props;

    return useCallback(
        (event: ChangeEvent<HTMLInputElement>) => {
            onOptionsChange({
                ...options,
                jsonData: {
                    ...options.jsonData,
                    [propertyName]: event.target.checked,
                },
            });
        },
        [onOptionsChange, options, propertyName]
    );
}
