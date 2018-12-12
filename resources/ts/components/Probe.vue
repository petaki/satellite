<template>
    <div class="flex-1 overflow-y-auto max-h-screen px-6 py-6">
        <div class="bg-indigo-dark px-6 py-6 rounded w-full max-w-md ml-auto mr-auto">
            <div class="flex items-center mb-6">
                <h4 class="text-indigo-lighter">
                    {{ isNew ? 'New' : 'Edit' }}
                </h4>
                <a v-for="(value, name, index) in types"
                   :key="value"
                   class="field-tab"
                   :class="{'ml-auto rounded-l': index === 0, 'rounded-r': index === typeSize - 1, active: probe.type === value}"
                   href="#"
                   @click.prevent="probe.type = value">
                    {{ name }}
                </a>
            </div>
            <form @submit.prevent="submit()">
                <div class="mb-6">
                    <label class="field-label required mb-2" for="name">
                        Name
                    </label>
                    <input id="name"
                           class="field"
                           type="text"
                           placeholder="Probe"
                           v-model="probe.name">
                    <p v-if="!isValid('name')" class="field-error mt-2">
                        {{ errors.name }}
                    </p>
                </div>
                <h4 class="text-indigo-lighter mb-6">
                    Redis
                </h4>
                <div class="mb-4">
                    <label class="field-label required mb-2" for="redis_host">
                        Redis Host
                    </label>
                    <input id="redis_host"
                           class="field"
                           type="text"
                           v-model="probe.redisHost">
                    <p v-if="!isValid('redisHost')" class="field-error mt-2">
                        {{ errors.redisHost }}
                    </p>
                </div>
                <div class="mb-4">
                    <label class="field-label required mb-2" for="redis_port">
                        Redis Port
                    </label>
                    <input id="redis_port"
                           class="field"
                           type="number"
                           placeholder="6379"
                           v-model.number="probe.redisPort">
                    <p v-if="!isValid('redisPort')" class="field-error mt-2">
                        {{ errors.redisPort }}
                    </p>
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="redis_password">
                        Redis Password
                    </label>
                    <input id="redis_password"
                           class="field"
                           type="password"
                           v-model="probe.redisPassword">
                    <p v-if="!isValid('redisPassword')" class="field-error mt-2">
                        {{ errors.redisPassword }}
                    </p>
                </div>
                <div class="mb-6">
                    <label class="field-label required mb-2" for="redis_key_prefix">
                        Redis Key Prefix
                    </label>
                    <input id="redis_key_prefix"
                           class="field"
                           type="text"
                           placeholder="probe:"
                           v-model="probe.redisKeyPrefix">
                    <p v-if="!isValid('redisKeyPrefix')" class="field-error mt-2">
                        {{ errors.redisKeyPrefix }}
                    </p>
                </div>
                <template v-if="isSSH">
                    <div class="flex items-center mb-6">
                        <h4 class="text-indigo-lighter">
                            SSH
                        </h4>
                        <a v-for="(value, name, index) in sshTypes"
                           :key="value"
                           class="field-tab"
                           :class="{'ml-auto rounded-l': index === 0, 'rounded-r': index === sshTypeSize - 1, active: probe.sshType === value}"
                           href="#"
                           @click.prevent="probe.sshType = value">
                            {{ name }}
                        </a>
                    </div>
                    <div class="mb-4">
                        <label class="field-label required mb-2" for="ssh_host">
                            SSH Host
                        </label>
                        <input id="ssh_host"
                               class="field"
                               type="text"
                               v-model="probe.sshHost">
                        <p v-if="!isValid('sshHost')" class="field-error mt-2">
                            {{ errors.sshHost }}
                        </p>
                    </div>
                    <div class="mb-4">
                        <label class="field-label required mb-2" for="ssh_port">
                            SSH Port
                        </label>
                        <input id="ssh_port"
                               class="field"
                               type="number"
                               placeholder="22"
                               v-model.number="probe.sshPort">
                        <p v-if="!isValid('sshPort')" class="field-error mt-2">
                            {{ errors.sshPort }}
                        </p>
                    </div>
                    <div class="mb-4">
                        <label class="field-label required mb-2" for="ssh_user">
                            SSH User
                        </label>
                        <input id="ssh_user"
                               class="field"
                               type="text"
                               v-model="probe.sshUser">
                        <p v-if="!isValid('sshUser')" class="field-error mt-2">
                            {{ errors.sshUser }}
                        </p>
                    </div>
                    <div v-if="isSSHPassword" class="mb-6">
                        <label class="field-label mb-2" for="ssh_password">
                            SSH Password
                        </label>
                        <input id="ssh_password"
                               class="field"
                               type="password"
                               v-model="probe.sshPassword">
                        <p v-if="!isValid('sshPassword')" class="field-error mt-2">
                            {{ errors.sshPassword }}
                        </p>
                    </div>
                    <template v-else>
                        <div class="mb-4">
                            <label class="field-label mb-2"
                                   for="ssh_key_file"
                                   @click="selectSSHKeyFile()">
                                SSH Key File
                            </label>
                            <input id="ssh_key_file"
                                   ref="ssh_key_file"
                                   class="hidden"
                                   type="file"
                                   @change="updateSSHKeyFile()">
                            <label class="block field"
                                   for="ssh_key_file"
                                   @click="selectSSHKeyFile()">
                                {{ probe.sshKeyFile || 'Choose a file...' }}
                            </label>
                            <p v-if="!isValid('sshKeyFile')" class="field-error mt-2">
                                {{ errors.sshKeyFile }}
                            </p>
                        </div>
                        <div class="mb-6">
                            <label class="field-label mb-2" for="ssh_key_passphrase">
                                SSH Key Passphrase
                            </label>
                            <input id="ssh_key_passphrase"
                                   class="field"
                                   type="password"
                                   v-model="probe.sshKeyPassphrase">
                            <p v-if="!isValid('sshKeyPassphrase')" class="field-error mt-2">
                                {{ errors.sshKeyPassphrase }}
                            </p>
                        </div>
                    </template>
                </template>
                <div class="flex">
                    <button class="field-btn" type="submit">
                        Save
                    </button>
                    <button v-if="!isNew"
                            class="field-remove-btn ml-auto"
                            type="button"
                            @click="removeProbe()">
                        Delete
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
    import _ from 'lodash';
    import { Component, Vue, Watch } from 'vue-property-decorator';
    import { IProbe, ISelected, ProbeType, SSHType } from '../store/types';
    import { Mutation, State } from 'vuex-class';

    @Component({
        name: 'Probe',
    })
    export default class Probe extends Vue {
        isNew: boolean = true;
        probe: IProbe = this.defaultProbe();
        errors: { [key: string]: string } = {};

        @State selected!: ISelected;

        @Mutation select!: (selected: ISelected) => void;
        @Mutation add!: (probe: IProbe) => void;
        @Mutation edit!: (probe: IProbe) => void;
        @Mutation remove!: (probe: IProbe) => void;

        get isSSH(): boolean {
            return this.probe.type === ProbeType.SSH;
        }

        get isSSHPassword(): boolean {
            return this.probe.sshType === SSHType.Password;
        }

        get types(): object {
            return _.pickBy(ProbeType, (value) => _.isNumber(value));
        }

        get typeSize(): number {
            return _.size(this.types);
        }

        get sshTypes(): object {
            return _.pickBy(SSHType, (value) => _.isNumber(value));
        }

        get sshTypeSize(): number {
            return _.size(this.sshTypes);
        }

        @Watch('selected')
        onSelect(selected: ISelected): void {
            this.isNew = !_.has(selected, 'probe');
            this.errors = {};

            if (!this.isNew) {
                this.probe = _.cloneDeep(selected.probe) as IProbe;
            } else {
                this.probe = this.defaultProbe();
            }
        }

        defaultProbe(): IProbe {
            return {
                name: 'Probe',
                type: ProbeType.Standard,
                redisHost: '',
                redisPort: 6379,
                redisPassword: '',
                redisKeyPrefix: 'probe:',
                sshType: SSHType.Password,
                sshHost: '',
                sshPort: 22,
                sshUser: '',
                sshPassword: '',
                sshKeyFile: '',
                sshKeyPassphrase: '',
            };
        }

        selectSSHKeyFile(): void {
            this.probe.sshKeyFile = '';
        }

        updateSSHKeyFile(): void {
            const files = _.get(this.$refs.ssh_key_file, 'files');

            this.probe.sshKeyFile = _.get(
                _.first(files), 'path'
            );
        }

        isValid(name?: string): boolean {
            if (name) {
                return !_.has(this.errors, name);
            }

            return _.isEmpty(this.errors);
        }

        validate(): void {
            this.errors = {};

            if (_.isEmpty(this.probe.name)) {
                this.errors.name = 'The Name field is required.';
            }

            if (_.isEmpty(this.probe.redisHost)) {
                this.errors.redisHost = 'The Redis Host field is required.';
            }

            if (!this.probe.redisPort) {
                this.errors.redisPort = 'The Redis Port field is required.';
            } else if (this.probe.redisPort < 1) {
                this.errors.redisPort = 'The Redis Port must be at least 1.';
            }

            if (_.isEmpty(this.probe.redisKeyPrefix)) {
                this.errors.redisKeyPrefix = 'The Redis Key Prefix field is required.';
            }

            if (this.isSSH) {
                if (_.isEmpty(this.probe.sshHost)) {
                    this.errors.sshHost = 'The SSH Host field is required.';
                }

                if (!this.probe.sshPort) {
                    this.errors.sshPort = 'The SSH Port field is required.';
                } else if (this.probe.sshPort < 1) {
                    this.errors.sshPort = 'The SSH Port must be at least 1.';
                }

                if (_.isEmpty(this.probe.sshUser)) {
                    this.errors.sshUser = 'The SSH User field is required.';
                }
            }
        }

        submit(): void {
            this.validate();

            if (this.isValid()) {
                const probe = _.cloneDeep(this.probe);

                if (this.isNew) {
                    this.add(probe);
                } else {
                    this.edit(probe);
                }

                this.select({
                    name: 'probe',
                    probe,
                });
            }
        }

        removeProbe(): void {
            this.remove(this.selected.probe as IProbe);
            this.select({
                name: 'new'
            });
        }
    };
</script>
