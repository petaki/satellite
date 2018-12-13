import _ from 'lodash';
import Vue from 'vue';
import Vuex from 'vuex';
import { actions } from './actions';
import { mutations } from './mutations';
import { plugins } from './plugins';
import { state } from './state';

Vue.use(Vuex);

export default new Vuex.Store({
    actions,
    mutations,
    plugins,
    state,
});
