<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len vue/attribute-hyphenation -->
        <div class="bg-white p-8 dark:bg-slate-700 dark:text-slate-300">
            <card-title>
                <cube-icon class="h-6 w-6" />
                <input v-model="keyword"
                       class="form-input flex-1 sm:mx-2"
                       placeholder="Search">
                <button class="btn-white"
                        type="button"
                        @click="router.reload()">
                    <arrow-path-icon class="h-6 w-6" />
                </button>
                <inertia-link class="btn-red sm:ml-2"
                              href="/probe/delete-all"
                              method="delete"
                              as="button"
                              :onBefore="confirmDelete">
                    <trash-icon class="h-6 w-6" />
                </inertia-link>
            </card-title>
            <div class="relative overflow-hidden overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr>
                            <th class="px-6 py-3 text-left">
                                Probe
                            </th>
                            <th class="px-6 py-3 text-left">
                                CPU
                            </th>
                            <th class="w-24 px-6 py-3 text-left"></th>
                            <th class="px-6 py-3 text-left">
                                Memory
                            </th>
                            <th class="w-24 px-6 py-3 text-left"></th>
                            <th class="px-6 py-3 text-left">
                                Load 1
                            </th>
                            <th class="px-6 py-3 text-left">
                                Load 5
                            </th>
                            <th class="px-6 py-3 text-left">
                                Load 15
                            </th>
                            <th class="px-6 py-3 text-left">
                                Status
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="!filteredProbes.length">
                            <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600"
                                colspan="9">
                                No probes.
                            </td>
                        </tr>
                        <tr v-for="probe in filteredProbes"
                            :key="probe.name"
                            class="cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800/50"
                            @click="router.visit(`/cpu?probe=${probe.name}`)">
                            <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                {{ probe.name }}
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                <div class="min-w-52">
                                    <div class="mb-2 flex items-center justify-between gap-3 text-sm">
                                        <span>{{ probe.cpu.toFixed(1) }}%</span>
                                    </div>
                                    <div class="h-2 rounded-full bg-slate-200 dark:bg-slate-800">
                                        <div class="h-2 rounded-full transition-all duration-1000 ease-out"
                                             :class="barColor(probe.cpu, probe.cpuAlarm, 'bg-emerald-500')"
                                             :style="{ width: animated ? `${clamp(probe.cpu)}%` : '0%' }"></div>
                                    </div>
                                </div>
                            </td>
                            <td class="w-24 border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                {{ probe.cpu.toFixed(1) }}%
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                <div class="min-w-52">
                                    <div class="mb-2 flex items-center justify-between gap-3 text-sm">
                                        <span>{{ probe.memory.toFixed(1) }}%</span>
                                    </div>
                                    <div class="h-2 rounded-full bg-slate-200 dark:bg-slate-800">
                                        <div class="h-2 rounded-full transition-all duration-1000 ease-out"
                                             :class="barColor(probe.memory, probe.memAlarm, 'bg-blue-500')"
                                             :style="{ width: animated ? `${clamp(probe.memory)}%` : '0%' }"></div>
                                    </div>
                                </div>
                            </td>
                            <td class="w-24 border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                {{ probe.memory.toFixed(1) }}%
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 tabular-nums dark:border-slate-600">
                                {{ probe.load1.toFixed(2) }}
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 tabular-nums dark:border-slate-600">
                                {{ probe.load5.toFixed(2) }}
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 tabular-nums dark:border-slate-600">
                                {{ probe.load15.toFixed(2) }}
                            </td>
                            <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600">
                                <span class="inline-flex h-3 w-3 rounded-full"
                                      :class="probe.hasBeat ? 'bg-emerald-500' : 'bg-red-500'">
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <!-- eslint-enable max-len vue/attribute-hyphenation -->
    </div>
</template>

<script setup lang="ts">
import { CubeIcon } from '@heroicons/vue/24/outline';
import { ArrowPathIcon, TrashIcon } from '@heroicons/vue/20/solid';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    watch
} from 'vue';

import { router } from '@inertiajs/vue3';
import type { ProbeSummary } from '../../types';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';

const {
    probes = []
} = defineProps<{
    probes?: ProbeSummary[]
}>();

defineOptions({
    layout: Layout
});

const keyword = ref('');
const subtitle = ref('Probes');
const reloadTimer = 60000;
const animated = ref(false);

let reloadInterval: ReturnType<typeof setInterval>;

const links = ref([
    { name: subtitle }
]);

const clamp = (value: number) => Math.min(Math.max(value, 0), 100);

const barColor = (value: number, alarm: number, color: string) => {
    if (alarm > 0) {
        if (value >= alarm) {
            return 'bg-red-500';
        }

        if (value >= alarm * 0.8) {
            return 'bg-amber-500';
        }
    }

    return color;
};

const filteredProbes = computed(() => {
    const words = keyword.value.trim().split(' ');

    return probes.filter(probe => {
        let has = true;

        words.forEach(word => {
            if (probe.name.includes(word)) {
                return true;
            }

            has = false;

            return false;
        });

        return has;
    });
});

const animateBars = () => {
    animated.value = false;
    requestAnimationFrame(() => { animated.value = true; });
};

// eslint-disable-next-line no-alert
const confirmDelete = () => window.confirm('Are you sure you want to delete all probes?');

onMounted(() => {
    animateBars();
    reloadInterval = setInterval(() => router.reload(), reloadTimer);
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});

watch(() => probes, animateBars);
</script>
