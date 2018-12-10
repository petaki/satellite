import _ from 'lodash';
import Vue from 'vue';
import Vuex, { MutationTree, Store } from 'vuex';
import { IProbe, ISelected, IState, STORAGE_KEY } from './types';

Vue.use(Vuex);

const state: IState = {
    selected: {
        name: 'new',
    },

    probes: JSON.parse(window.localStorage.getItem(STORAGE_KEY) || '[]'),
};

const mutations: MutationTree<IState> = {
    select(current: IState, selected: ISelected) {
        current.selected = selected;
    },

    add(current: IState, probe: IProbe) {
        current.probes.push(probe);
    },

    edit(current: IState, probe: IProbe) {
        current.probes.splice(current.probes.indexOf(probe), 1, probe);
    },

    remove(current: IState, probe: IProbe) {
        current.probes.splice(current.probes.indexOf(probe), 1);
    },
};

const localStoragePlugin = (store: Store<IState>) => {
    store.subscribe((mutation, { probes }) => {
        window.localStorage.setItem(STORAGE_KEY, JSON.stringify(probes));
    });
};

const plugins = [
    localStoragePlugin,
];

export default new Vuex.Store({
    mutations,
    plugins,
    state,
});
