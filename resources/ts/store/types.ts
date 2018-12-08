export interface IState {
    probes: IProbe[];
}

export interface IProbe {
    name: string;
    type: ProbeType;
    sshType: SSHType;
}

export enum SSHType {
    Password,
    Key,
}

export enum ProbeType {
    Standard,
    SSH,
}
