import { ChartType } from '../store/types';

export interface IRepository {
    findCpuDataset(chartType: ChartType): Promise<IDataset>;
    findMemoryDataset(chartType: ChartType): Promise<IDataset>;
    findDiskDataset(chartType: ChartType, path: string): Promise<IDataset>;
}

export interface IDataset {
    [key: string]: string;
}
