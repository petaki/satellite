import _ from 'lodash';
import Vue from 'vue';
import Vuex, { MutationTree } from 'vuex';
import { IProbe, ISelected, IState } from './types';

Vue.use(Vuex);

const state: IState = {
    selected: {
        name: 'new',
    },

    probes: [],
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

export default new Vuex.Store({
    mutations,
    state,
});
