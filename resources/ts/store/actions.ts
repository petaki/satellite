import _ from 'lodash';
import { ActionContext, ActionTree } from 'vuex';
import { IFlash, IState } from './types';

export const actions: ActionTree<IState, any> = {
    flash(context: ActionContext<IState, any>, flash: IFlash) {
        context.commit('flash', flash);

        if (!_.isUndefined(flash.timeout)) {
            setTimeout(() => context.commit('flash'), flash.timeout);
        }
    },
};
