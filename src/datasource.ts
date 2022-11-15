import { DataSourceInstanceSettings } from '@grafana/data';
import { NumaflowDataQuery, NumaflowDataSourceOptions } from './types';
import { DataSourceWithBackend } from "@grafana/runtime";

export class NumaflowDataSource extends DataSourceWithBackend<NumaflowDataQuery, NumaflowDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NumaflowDataSourceOptions>) {
    super(instanceSettings);
  }

  // async query(options: DataQueryRequest<NumaflowDataQuery>): Promise<DataQueryResponse> {
  //   const { range } = options;
  //   const from = range!.from.valueOf();
  //   const to = range!.to.valueOf();
  //
  //   // Return a constant for each query.
  //   const data = options.targets.map((target) => {
  //     const query = defaults(target, defaultQuery);
  //     return new MutableDataFrame({
  //       refId: query.refId,
  //       fields: [
  //         { name: 'Time', values: [from, to], type: FieldType.time },
  //         { name: 'Value', values: [query.constant, query.constant], type: FieldType.number },
  //       ],
  //     });
  //   });
  //
  //   return { data };
  // }
}
