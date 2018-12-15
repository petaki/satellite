import { RedisClient } from 'redis';

export const STORAGE_KEY = '_probes';

export interface IState {
    chartType: ChartType;
    connection?: IConnection;
    flash?: IFlash;
    selected: ISelected;
    probes: IProbe[];
}

export interface IConnection {
    probe: IProbe;
    client: RedisClient;
}

export interface ISelected {
    name: string;
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
