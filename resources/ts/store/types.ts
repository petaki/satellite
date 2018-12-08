export interface IState {
    probes: IProbe[];
}

export interface IProbe {
    isConnected: boolean;
    name: string;
}
