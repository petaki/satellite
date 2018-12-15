import _ from 'lodash';
import moment, { Moment } from 'moment';
import { RedisClient } from 'redis';
import { promisify } from 'util';
import { ChartType, IDataset, IProbe, IRepository } from '../store/types';

export default class Repository implements IRepository {
    protected probe: IProbe;
    protected client: RedisClient;

    constructor(probe: IProbe, client: RedisClient) {
        this.probe = probe;
        this.client = client;
    }

    public async findCpuDataset(chartType: ChartType): Promise<IDataset> {
        const batch = this.client.batch();
        const execAsync = promisify(batch.exec).bind(batch);

        _.forEach(
            this.timestamps(chartType), (timestamp) => batch.hgetall(`cpu:${timestamp}`),
        );

        const data = await execAsync();

        return _.assignIn({}, ...data);
    }

    public async findMemoryDataset(chartType: ChartType): Promise<IDataset> {
        const batch = this.client.batch();
        const execAsync = promisify(batch.exec).bind(batch);

        _.forEach(
            this.timestamps(chartType), (timestamp) => batch.hgetall(`memory:${timestamp}`),
        );

        const data = await execAsync();

        return _.assignIn({}, ...data);
    }

    public async findDiskPaths(cursor: string = '0', paths: string[] = []): Promise<string[]> {
        const prefix = `${this.probe.redisKeyPrefix}disk:${this.today().unix()}:`;
        const scanAsync = promisify(this.client.scan).bind(this.client);

        const data = await scanAsync(
            cursor,
            'MATCH',
            `${prefix}*`,
            'COUNT',
            '10',
        );

        if (data[0] !== '0') {
            const newPaths = _.map<string, string>(
                data[1], (path) => Buffer.from(path.replace(prefix, ''), 'base64').toString(),
            );

            return await this.findDiskPaths(
                data[0], [...paths, ...newPaths],
            );
        }

        return paths.sort();
    }

    public async findDiskDataset(chartType: ChartType, path: string): Promise<IDataset> {
        const prefix = Buffer.from(path).toString('base64');
        const batch = this.client.batch();
        const execAsync = promisify(batch.exec).bind(batch);

        _.forEach(
            this.timestamps(chartType), (timestamp) => batch.hgetall(`disk:${timestamp}:${prefix}`),
        );

        const data = await execAsync();

        return _.assignIn({}, ...data);
    }

    protected today(): Moment {
        return moment().startOf('day');
    }

    protected timestamps(chartType: ChartType): number[] {
        const timestamps = [];

        const end = this.today();
        let start = moment(end);

        if (chartType === ChartType.Week) {
            start = moment(end).subtract(6, 'weeks');
        } else if (chartType === ChartType.Month) {
            start = moment(end).subtract(1, 'months');
        }

        for (const current = moment(start); current.diff(end, 'days') <= 0; current.add(1, 'days')) {
            timestamps.push(current.unix());
        }

        return timestamps;
    }
}
