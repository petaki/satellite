import { MutationTree } from 'vuex';
import { IConnection, IFlash, IProbe, ISelected, IState } from './types';

export const mutations: MutationTree<IState> = {
    connection(state: IState, connection?: IConnection) {
        state.connection = connection;
    },

    select(state: IState, selected: ISelected) {
        state.selected = selected;
    },

    flash(state: IState, flash?: IFlash) {
        state.flash = flash;
    },

    add(state: IState, probe: IProbe) {
        state.probes.push(probe);
    },

    edit(state: IState, { newProbe, oldProbe }) {
        state.probes.splice(state.probes.indexOf(oldProbe), 1, newProbe);
    },

    remove(state: IState, probe: IProbe) {
        state.probes.splice(state.probes.indexOf(probe), 1);
    },
};
