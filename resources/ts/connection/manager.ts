import _ from 'lodash';
import { IConnection, IManager } from './types';

export default class Manager implements IManager {
    protected connection: IConnection | undefined;

    get isConnected(): boolean {
        return !_.isEmpty(this.connection);
    }
}
