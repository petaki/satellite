import _ from 'lodash';
import Vue from 'vue';
import Vuex from 'vuex';
import { IState } from './types';

Vue.use(Vuex);

const state: IState = {
    probes: [],
};

export default new Vuex.Store({
    state,
});
