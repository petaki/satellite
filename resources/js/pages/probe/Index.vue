<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len vue/attribute-hyphenation -->
        <div class="flex items-center bg-white p-5 shadow-sm mb-5">
            <div class="flex-1">
                <input v-model="keyword"
                       class="bg-transparent border-slate-300 text-slate-600 placeholder-gray-400 rounded-sm focus:border-cyan-500 focus:ring-cyan-500 w-full h-11"
                       type="text"
                       placeholder="Search">
            </div>
            <button class="btn-white px-3 py-0 h-11 mx-2"
                    type="button"
                    @click="router.reload()">
                <arrow-path-icon class="h-5 w-5" />
            </button>
            <inertia-link class="btn-red px-3 py-0 h-11"
                          href="/probe/delete-all"
                          method="delete"
                          as="button"
                          :onBefore="confirmDelete">
                <trash-icon class="h-5 w-5" />
            </inertia-link>
        </div>
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-if="!probes.length"
                 class="bg-white p-8">
                No probes.
            </div>
            <inertia-link v-for="probe in filteredProbes"
                          :key="probe"
                          :href="`/cpu?probe=${probe}`"
                          class="bg-white p-8 flex flex-col sm:flex-row items-center text-base text-gray-800 text-lg hover:text-cyan-500">
                <cube-icon class="h-6 w-6 sm:mr-2" />
                <span class="break-all flex-1 sm:mr-auto">
                    {{ probe }}
                </span>
                <chevron-right-icon class="h-6 w-6 sm:ml-2" />
            </inertia-link>
        </div>
        <!-- eslint-enable max-len vue/attribute-hyphenation -->
    </div>
</template>

<script setup>
import {
    CubeIcon,
    ChevronRightIcon
} from '@heroicons/vue/24/outline';

import { ArrowPathIcon, TrashIcon } from '@heroicons/vue/20/solid';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    defineProps,
    defineOptions
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import Layout from '../../base/Layout.vue';

const { probes } = defineProps({
    probes: {
        type: Array,
        default: () => []
    }
});

defineOptions({
    layout: Layout
});

const keyword = ref('');
const subtitle = ref('Probes');
const reloadTimer = 60000;

let reloadInterval;

const links = ref([
    { name: subtitle }
]);

const filteredProbes = computed(() => {
    const words = keyword.value.trim().split(' ');

    return probes.filter(probe => {
        let has = true;

        words.forEach(word => {
            if (probe.includes(word)) {
                return true;
            }

            has = false;

            return false;
        });

        return has;
    });
});

// eslint-disable-next-line no-alert
const confirmDelete = () => window.confirm('Are you sure you want to delete all probes?');

onMounted(() => {
    reloadInterval = setInterval(() => router.reload(), reloadTimer);
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});
</script>
