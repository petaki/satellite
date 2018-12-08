import _ from 'lodash';
import Vue from 'vue';
import Vuex, { GetterTree } from 'vuex';
import { IProbe, IState } from './types';

Vue.use(Vuex);

const state: IState = {
    probes: [],
};

const getters: GetterTree<IState, any> = {
    connected: (s) => _.find<IProbe>(s.probes, (probe) => probe.isConnected),
};

export default new Vuex.Store({
    getters,
    state,
});
