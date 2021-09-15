<template>
    <inertia-head :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <chip-icon class="h-6 w-6 mr-2" />
                <span class="flex-1 mr-auto">
                    {{ subtitle }}
                </span>
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? '/'
                                  : `/?type=${type.value}`"
                              class="hover:text-cyan-500 ml-3"
                              :class="{'text-cyan-500': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
            </card-title>
            <div class="chart">
                <apexchart v-if="cpuSeries"
                           type="line"
                           :series="series"
                           height="100%"
                           :options="options" />
            </div>
        </div>
    </div>
</template>

<script>
import {
    ChipIcon
} from '@heroicons/vue/outline';

import { ref, toRefs, computed, onMounted, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        ChipIcon,
        Breadcrumb,
        CardTitle
    },

    props: {
        seriesType: {
            type: String,
            default: ''
        },

        seriesTypes: {
            type: Array,
            default: () => []
        },

        cpuSeries: {
            type: Array,
            default: () => []
        }
    },

    layout: Layout,

    setup(props) {
        const { cpuSeries } = toRefs(props);
        const subtitle = ref('CPU');
        const reloadInterval = ref();
        const reloadTimer = ref(60000);
        const options = ref({});

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [{
            name: 'CPU',
            data: cpuSeries.value
        }]);

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
            options,
            links,
            series
        };
    }
};
</script>
