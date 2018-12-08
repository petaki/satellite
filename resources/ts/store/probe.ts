import { IProbe, ProbeType, SSHType } from './types';

export default class Probe implements IProbe {
    public name: string;
    public type: ProbeType = ProbeType.Standard;
    public sshType: SSHType = SSHType.Password;

    protected _isConnected: boolean = false;

    get isConnected(): boolean {
        return this._isConnected;
    }
}
