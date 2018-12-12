import { Store } from 'vuex';
import { IState, STORAGE_KEY } from './types';

const localStoragePlugin = (store: Store<IState>) => {
    store.subscribe((mutation, { probes }) => {
        window.localStorage.setItem(STORAGE_KEY, JSON.stringify(probes));
    });
};

export const plugins = [
    localStoragePlugin,
];
