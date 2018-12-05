<template>
    <div class="flex-1 px-6 py-6">
        <div class="bg-indigo-dark px-6 py-6 rounded w-full max-w-sm">
            <h4 class="text-indigo-lighter mb-6">
                New
            </h4>
            <form>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="name">
                        Name
                    </label>
                    <input id="name"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="redis_host">
                        Redis Host
                    </label>
                    <input id="redis_host"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text"
                           placeholder="localhost">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="redis_port">
                        Redis Port
                    </label>
                    <input id="redis_port"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text"
                           placeholder="6379">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="redis_password">
                        Redis Password
                    </label>
                    <input id="redis_password"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="password">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_host">
                        SSH Host
                    </label>
                    <input id="ssh_host"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_port">
                        SSH Port
                    </label>
                    <input id="ssh_port"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text"
                           placeholder="22">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_user">
                        SSH User
                    </label>
                    <input id="ssh_user"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_password">
                        SSH Password
                    </label>
                    <input id="ssh_password"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="password">
                </div>
                <div class="mb-4">
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_key_file">
                        SSH Key File
                    </label>
                    <input id="ssh_key_file"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="text">
                </div>
                <div>
                    <label class="block text-indigo-light text-sm mb-2" for="ssh_key_passphrase">
                        SSH Key Passphrase
                    </label>
                    <input id="ssh_key_passphrase"
                           class="appearance-none leading-tight text-indigo bg-indigo-darkest rounded px-4 py-2 w-full focus:outline-none"
                           type="password">
                </div>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
    import { Component, Vue } from "vue-property-decorator";

    @Component({
        name: "Probe",
    })
    export default class Probe extends Vue {};
</script>
