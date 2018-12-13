import { IState, STORAGE_KEY } from './types';

export const state: IState = {
    connection: undefined,
    flash: undefined,
    probes: JSON.parse(window.localStorage.getItem(STORAGE_KEY) || '[]'),

    selected: {
        name: 'new',
    },
};
