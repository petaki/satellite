<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len -->
        <div class="mb-5 xl:w-1/4 xl:pr-4">
            <input v-model="keyword"
                   class="bg-transparent border-slate-300 text-slate-600 placeholder-gray-400 rounded-sm focus:border-cyan-500 focus:ring-cyan-500 w-full"
                   type="text"
                   placeholder="Search">
        </div>
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-if="!probes.length"
                 class="bg-white p-8">
                No probes.
            </div>
            <a v-for="probe in filteredProbes"
               :key="probe"
               :href="`/cpu?probe=${probe}`"
               class="bg-white p-8 flex flex-col sm:flex-row items-center text-base text-gray-800 text-lg hover:text-cyan-500">
                <cube-icon class="h-6 w-6 sm:mr-2" />
                <span class="break-all flex-1 sm:mr-auto">
                    {{ probe }}
                </span>
                <chevron-right-icon class="h-6 w-6 sm:ml-2" />
            </a>
        </div>
        <!-- eslint-enable max-len -->
    </div>
</template>

<script>
import {
    CubeIcon,
    ChevronRightIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    toRefs,
    computed,
    onMounted,
    onUnmounted
} from 'vue';

import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        CubeIcon,
        ChevronRightIcon,
        Breadcrumb
    },

    layout: Layout,

    props: {
        probes: {
            type: Array,
            default: () => []
        }
    },

    setup(props) {
        const { probes } = toRefs(props);
        const keyword = ref('');
        const subtitle = ref('Probes');
        const reloadTimer = 60000;

        let reloadInterval;

        const links = ref([
            { name: subtitle }
        ]);

        const filteredProbes = computed(() => {
            const words = keyword.value.trim().split(' ');

            return probes.value.filter(probe => {
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

        onMounted(() => {
            reloadInterval = setInterval(() => Inertia.reload(), reloadTimer);
        });

        onUnmounted(() => {
            clearInterval(reloadInterval);
        });

        return {
            subtitle,
            links,
            keyword,
            filteredProbes
        };
    }
};
</script>
