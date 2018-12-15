import _ from 'lodash';
import moment from 'moment';
import { promisify } from 'util';
import { ChartType, IConnection } from '../store/types';
import { IDataset, IRepository } from './types';

export default class Repository implements IRepository {
    protected connection?: IConnection;

    constructor(connection?: IConnection) {
        this.connection = connection;
    }

    public async findCpuDataset(chartType: ChartType): Promise<IDataset> {
        if (_.isUndefined(this.connection)) {
            return {};
        }

        const batch = this.connection.client.batch();
        const execAsync = promisify(batch.exec).bind(batch);

        _.forEach(
            this.timestamps(chartType), (timestamp) => batch.hgetall(`cpu:${timestamp}`),
        );

        const data = await execAsync();

        return _.assignIn({}, ...data);
    }

    public async findMemoryDataset(chartType: ChartType): Promise<IDataset> {
        if (_.isUndefined(this.connection)) {
            return {};
        }

        const batch = this.connection.client.batch();
        const execAsync = promisify(batch.exec).bind(batch);

        _.forEach(
            this.timestamps(chartType), (timestamp) => batch.hgetall(`memory:${timestamp}`),
        );

        const data = await execAsync();

        return _.assignIn({}, ...data);
    }

    public async findDiskDataset(chartType: ChartType, path: string): Promise<IDataset> {
        return {};
    }

    protected timestamps(chartType: ChartType): number[] {
        const timestamps = [];

        const end = moment().startOf('day');
        let start = moment(end);

        if (chartType === ChartType.Week) {
            start = moment(end).subtract(6, 'days');
        } else if (chartType === ChartType.Month) {
            start = moment(end).subtract(1, 'months');
        }

        for (const current = moment(start); current.diff(end, 'days') <= 0; current.add(1, 'days')) {
            timestamps.push(current.unix());
        }

        return timestamps;
    }
}
