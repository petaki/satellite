export interface IState {
    selected: ISelected;
    probes: IProbe[];
}

export interface ISelected {
    name: string;
    probe?: IProbe;
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

export enum SSHType {
    Password,
    Key,
}

export enum ProbeType {
    Standard,
    SSH,
}
