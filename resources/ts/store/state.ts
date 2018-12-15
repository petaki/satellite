import { ChartType, IState, STORAGE_KEY } from './types';

export const state: IState = {
    chartType: ChartType.Day,
    connection: undefined,
    flash: undefined,
    probes: JSON.parse(window.localStorage.getItem(STORAGE_KEY) || '[]'),

    selected: {
        name: 'New',
    },
};
