<template>
    <div class="flex-1 overflow-y-auto max-h-screen px-6 py-6">
        <div class="bg-indigo-dark px-6 py-6 rounded w-full max-w-md ml-auto mr-auto">
            <div class="flex items-center mb-6">
                <h4 class="text-indigo-lighter">
                    New
                </h4>
                <a class="field-tab left active ml-auto" href="#">
                    Standard
                </a>
                <a class="field-tab right" href="#">
                    SSH
                </a>
            </div>
            <form>
                <div class="mb-6">
                    <label class="field-label required mb-2" for="name">
                        Name
                    </label>
                    <input id="name"
                           class="field"
                           type="text"
                           placeholder="Probe">
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
                           placeholder="localhost">
                </div>
                <div class="mb-4">
                    <label class="field-label required mb-2" for="redis_port">
                        Redis Port
                    </label>
                    <input id="redis_port"
                           class="field"
                           type="text"
                           placeholder="6379">
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="redis_password">
                        Redis Password
                    </label>
                    <input id="redis_password"
                           class="field"
                           type="password">
                </div>
                <div class="mb-6">
                    <label class="field-label required mb-2" for="redis_key_prefix">
                        Redis Key Prefix
                    </label>
                    <input id="redis_key_prefix"
                           class="field"
                           type="text"
                           placeholder="probe:">
                </div>
                <div class="flex items-center mb-6">
                    <h4 class="text-indigo-lighter">
                        SSH
                    </h4>
                    <a class="field-tab left active ml-auto" href="#">
                        Password
                    </a>
                    <a class="field-tab right" href="#">
                        Key
                    </a>
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="ssh_host">
                        SSH Host
                    </label>
                    <input id="ssh_host"
                           class="field"
                           type="text">
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="ssh_port">
                        SSH Port
                    </label>
                    <input id="ssh_port"
                           class="field"
                           type="text"
                           placeholder="22">
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="ssh_user">
                        SSH User
                    </label>
                    <input id="ssh_user"
                           class="field"
                           type="text">
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="ssh_password">
                        SSH Password
                    </label>
                    <input id="ssh_password"
                           class="field"
                           type="password">
                </div>
                <div class="mb-4">
                    <label class="field-label mb-2" for="ssh_key_file">
                        SSH Key File
                    </label>
                    <input id="ssh_key_file"
                           class="field"
                           type="text">
                </div>
                <div class="mb-6">
                    <label class="field-label mb-2" for="ssh_key_passphrase">
                        SSH Key Passphrase
                    </label>
                    <input id="ssh_key_passphrase"
                           class="field"
                           type="password">
                </div>
                <div class="text-center">
                    <button class="field-btn" type="submit">
                        Save and Connect
                    </button>
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
