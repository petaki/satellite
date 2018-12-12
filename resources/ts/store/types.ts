export const STORAGE_KEY = '_probes';

export interface IState {
    flash?: IFlash;
    selected: ISelected;
    probes: IProbe[];
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
    name: string;
    type: ProbeType;
    redisHost: string;
    redisPort: number;
    redisPassword: string;
    redisKeyPrefix: string;
    sshType: SSHType;
    sshHost: string;
    sshPort: number;
    sshUser: string;
    sshPassword: string;
    sshKeyFile: string;
    sshKeyPassphrase: string;
}

export enum FlashType {
    Success,
    Loading,
    Error,
}

export enum SSHType {
    Password,
    Key,
}

export enum ProbeType {
    Standard,
    SSH,
}
