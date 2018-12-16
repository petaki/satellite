import fs from 'fs';
import _ from 'lodash';
import net, { AddressInfo, Server } from 'net';
import redis, { ClientOpts, RedisClient } from 'redis';
import { Client, ConnectConfig } from 'ssh2';
import { ActionContext, ActionTree } from 'vuex';
import Repository from '../support/repository';
import { FlashType, IFlash, IProbe, IState, ProbeType, RedisType, SSHType } from './types';

const redisCreateOptions = (probe: IProbe): ClientOpts => {
    const options: ClientOpts = {
        db: probe.redisDatabase,
        host: probe.redisHost,
        port: probe.redisPort,
        prefix: probe.redisKeyPrefix,
    };

    if (!_.isEmpty(probe.redisPassword)) {
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

    return options;
};

const redisCreate = async (redisOptions: ClientOpts): Promise<RedisClient> => {
    return new Promise<RedisClient>((resolve, reject) => {
        const client = redis.createClient(redisOptions);

        client.on('error', reject);
        client.on('connect', () => resolve(client));
    });
};

const sshCreateOptions = (probe: IProbe): ConnectConfig => {
    const options: ConnectConfig = {
        host: probe.sshHost,
        port: probe.sshPort,
        username: probe.sshUser,
    };

    if (probe.sshType === SSHType.Key) {
        options.privateKey = fs.readFileSync(probe.sshKeyFile);

        if (!_.isEmpty(probe.sshKeyPassphrase)) {
            options.passphrase = probe.sshKeyPassphrase;
        }
    } else if (!_.isEmpty(probe.sshPassword)) {
        options.password = probe.sshPassword;
    }

    return options;
};

const sshConnect = async (probe: IProbe): Promise<Client> => {
    return new Promise<Client>((resolve, reject) => {
        const sshClient = new Client();

        sshClient
            .on('ready', () => resolve(sshClient))
            .on('error', reject)
            .connect(sshCreateOptions(probe));
    });
};

const serverCreate = async (redisOptions: ClientOpts, sshClient: Client): Promise<Server> => {
    return new Promise<Server>((resolve, reject) => {
        const server = net.createServer((sock) => {
            sshClient.forwardOut(
                sock.remoteAddress,
                sock.remotePort,
                redisOptions.host,
                redisOptions.port,
                (err, stream) => {
                    if (err) {
                        sock.end();
                    } else {
                        sock.pipe(stream).pipe(sock);
                    }
                },
            );
        });

        server.on('error', reject).listen(0, () => resolve(server));
    });
};

export const actions: ActionTree<IState, any> = {
    async connect(context: ActionContext<IState, any>, probe: IProbe) {
        context.dispatch('flash', {
            message: 'Connecting...',
            type: FlashType.Loading,
        });

        const redisOptions = redisCreateOptions(probe);

        try {
            let client: RedisClient;
            let sshClient: Client;
            let server: Server;

            if (probe.type === ProbeType.SSH) {
                sshClient = await sshConnect(probe);
                server = await serverCreate(redisOptions, sshClient);

                client = await redisCreate({
                    ...redisOptions,
                    host: '127.0.0.1',
                    port: (server.address() as AddressInfo).port,
                });
            } else {
                client = await redisCreate(redisOptions);
            }

            const repository = new Repository(probe, client);

            context.commit('connection', {
                client,
                probe,
                repository,
                server,
                sshClient,
            });

            context.commit('select', {
                name: 'CPU',
            });

            context.dispatch('flash', {
                message: 'Connected',
                timeout: 1500,
                type: FlashType.Success,
            });
        } catch (error) {
            context.dispatch('flash', {
                message: `${error}`,
                timeout: 1500,
                type: FlashType.Error,
            });

            context.dispatch('disconnect');
        }
    },

    disconnect(context: ActionContext<IState, any>): void {
        if (_.isUndefined(context.state.connection)) {
            return;
        }

        context.state.connection.client.end(true);

        if (!_.isUndefined(context.state.connection.server)) {
            context.state.connection.server.close();
        }

        if (!_.isUndefined(context.state.connection.sshClient)) {
            context.state.connection.sshClient.end();
        }

        context.commit('connection');
        context.commit('select', {
            name: 'New',
        });
    },

    flash(context: ActionContext<IState, any>, flash: IFlash): void {
        context.commit('flash', flash);

        if (_.has(flash, 'timeout')) {
            setTimeout(() => context.commit('flash'), flash.timeout);
        }
    },
};
