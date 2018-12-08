import { Redis } from 'ioredis';
import { IProbe } from '../store/types';

export interface IManager {
    isConnected: boolean;
}

export interface IConnection {
    probe: IProbe;
    redis: Redis;
}
