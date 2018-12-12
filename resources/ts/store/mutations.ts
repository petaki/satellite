import { MutationTree } from 'vuex';
import { IFlash, IProbe, ISelected, IState } from './types';

export const mutations: MutationTree<IState> = {
    select(state: IState, selected: ISelected) {
        state.selected = selected;
    },

    flash(state: IState, flash?: IFlash) {
        state.flash = flash;
    },

    add(state: IState, probe: IProbe) {
        state.probes.push(probe);
    },

    edit(state: IState, probe: IProbe) {
        state.probes.splice(state.probes.indexOf(probe), 1, probe);
    },

    remove(state: IState, probe: IProbe) {
        state.probes.splice(state.probes.indexOf(probe), 1);
    },
};
