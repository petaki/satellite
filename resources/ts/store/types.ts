export interface IState {
    probes: IProbe[];
}

export enum ProbeType {
    Standard,
    SSH,
}

export enum SSHType {
    Password,
    Key,
}

export interface IProbe {
    isConnected: boolean;
    name: string;
    type: ProbeType;
    sshType: SSHType;
}
