import { Server } from 'net';
import { RedisClient } from 'redis';
import { Client } from 'ssh2';

export const STORAGE_KEY = '_probes';

export interface IState {
    chartType: ChartType;
    connection?: IConnection;
    flash?: IFlash;
    selected: ISelected;
    probes: IProbe[];
}

export interface IConnection {
    client: RedisClient;
    probe: IProbe;
    repository: IRepository;
    server?: Server;
    sshClient?: Client;
}

export interface IRepository {
    findCpuDataset(chartType: ChartType): Promise<IDataset>;
    findMemoryDataset(chartType: ChartType): Promise<IDataset>;
    findDiskPaths(cursor?: string, data?: string[]): Promise<string[]>;
    findDiskDataset(chartType: ChartType, path: string): Promise<IDataset>;
}

export interface IDataset {
    [key: string]: string;
}

export interface ISelected {
    name: string;
    path?: string;
    probe?: IProbe;
}

export interface IFlash {
    type: FlashType;
    timeout?: number;
    message: string;
}

export interface IProbe {
    type: ProbeType;
    name: string;
    redisType: RedisType;
    redisHost: string;
    redisPort: number;
    redisDatabase: number;
    redisPassword: string;
    redisKeyPrefix: string;
    redisKeyFile: string;
    redisCertificate: string;
    redisCaCert: string;
    sshType: SSHType;
    sshHost: string;
    sshPort: number;
    sshUser: string;
    sshPassword: string;
    sshKeyFile: string;
    sshKeyPassphrase: string;
}

export enum ChartType {
    Day,
    Week,
    Month,
}

export enum FlashType {
    Success,
    Loading,
    Error,
}

export enum ProbeType {
    Standard,
    SSH,
}

export enum RedisType {
    Normal,
    SSL,
}

export enum SSHType {
    Password,
    Key,
}
