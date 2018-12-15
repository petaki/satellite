import fs from 'fs';
import _ from 'lodash';
import redis, { ClientOpts } from 'redis';
import { ActionContext, ActionTree } from 'vuex';
import { FlashType, IFlash, IProbe, IState, RedisType } from './types';

export const actions: ActionTree<IState, any> = {
    connect(context: ActionContext<IState, any>, probe: IProbe) {
        context.dispatch('flash', {
            message: 'Connecting...',
            type: FlashType.Loading,
        });

        const options: ClientOpts = {
            db: probe.redisDatabase,
            host: probe.redisHost,
            port: probe.redisPort,
            prefix: probe.redisKeyPrefix,
        };

        if (probe.redisPassword) {
            options.password = probe.redisPassword;
        }

        if (probe.redisType === RedisType.SSL) {
            options.tls = {
                cert: fs.readFileSync(probe.redisCertificate),
                key: fs.readFileSync(probe.redisKeyFile),
            };

            if (!_.isEmpty(probe.redisCaCert)) {
                options.tls.ca = [
                    fs.readFileSync(probe.redisCaCert),
                ];
            }
        }

        const client = redis.createClient(options);

        client.once('error', (error) => {
            context.dispatch('flash', {
                message: `${error}`,
                timeout: 1500,
                type: FlashType.Error,
            });

            client.end();
        });

        client.once('connect', () => {
            context.commit('connection', {
                client,
                probe,
            });

            context.commit('select', {
                name: 'CPU',
            });

            context.dispatch('flash', {
                message: 'Connected',
                timeout: 1500,
                type: FlashType.Success,
            });
        });

        client.once('end', () => {
            context.commit('connection');
        });
    },

    flash(context: ActionContext<IState, any>, flash: IFlash) {
        context.commit('flash', flash);

        if (_.has(flash, 'timeout')) {
            setTimeout(() => context.commit('flash'), flash.timeout);
        }
    },
};
