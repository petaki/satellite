<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len -->
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-if="!probes.length"
                 class="bg-white p-8">
                No probes.
            </div>
            <a v-for="probe in probes"
               :key="probe"
               :href="`/cpu?probe=${probe}`"
               class="bg-white p-8 flex flex-col sm:flex-row items-center text-base text-gray-800 text-lg hover:text-cyan-500">
                <cube-icon class="h-6 w-6 sm:mr-2" />
                <span class="flex-1 sm:mr-auto">
                    {{ probe }}
                </span>
            </a>
        </div>
        <!-- eslint-enable max-len -->
    </div>
</template>

<script>
import {
    CubeIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    onMounted,
    onUnmounted
} from 'vue';

import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        CubeIcon,
        Breadcrumb
    },

    layout: Layout,

    props: {
        probes: {
            type: Array,
            default: () => []
        }
    },

    setup() {
        const subtitle = ref('Probes');
        const reloadInterval = ref();
        const reloadTimer = ref(60000);

        const links = ref([
            { name: subtitle }
        ]);

        onMounted(() => {
            reloadInterval.value = setInterval(() => Inertia.reload(), reloadTimer.value);
        });

        onUnmounted(() => {
            clearInterval(reloadInterval.value);
        });

        return {
            subtitle,
            reloadInterval,
            reloadTimer,
            links
        };
    }
};
</script>
