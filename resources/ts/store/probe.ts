import { IProbe } from './types';

export default class Probe implements IProbe {
    public name: string;
    protected _isConnected: boolean = false;

    get isConnected(): boolean {
        return this._isConnected;
    }
}
