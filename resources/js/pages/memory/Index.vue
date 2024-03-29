<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <document-duplicate-icon class="h-6 w-6 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
                <span class="text-slate-600 sm:mx-auto">
                    Point Interval:
                    <span class="text-cyan-500">
                        {{ duration(chunkSize) }}
                    </span>
                </span>
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? `/memory?probe=${probe}`
                                  : `/memory?probe=${probe}&type=${type.value}`"
                              class="hover:text-cyan-500 sm:ml-3"
                              :class="{'text-cyan-500': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
            </card-title>
            <div class="chart">
                <apexchart v-if="memoryMinSeries"
                           ref="chartEl"
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
    DocumentDuplicateIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    toRefs,
    computed,
    onMounted,
    onUnmounted, nextTick
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';
import useAnnotation from '../../common/useAnnotation';

export default {
    components: {
        DocumentDuplicateIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        probe: {
            type: String,
            required: true
        },

        chunkSize: {
            type: Number,
            required: true
        },

        seriesType: {
            type: String,
            default: ''
        },

        seriesTypes: {
            type: Array,
            default: () => []
        },

        memoryMinSeries: {
            type: Array,
            default: () => []
        },

        memoryMaxSeries: {
            type: Array,
            default: () => []
        },

        memoryAvgSeries: {
            type: Array,
            default: () => []
        },

        memoryAlarm: {
            type: Number,
            default: 0
        }
    },

    setup(props) {
        const {
            memoryMinSeries,
            memoryMaxSeries,
            memoryAvgSeries,
            memoryAlarm
        } = toRefs(props);

        const { alarm, max } = useAnnotation();
        const subtitle = ref('Memory');
        const chartEl = ref();
        const options = ref({});
        const reloadTimer = 60000;
        let reloadInterval;

        options.value.annotations = {
            yaxis: [
                max(Math.max(...(memoryMaxSeries.value ?? []).map(value => value.y)))
            ]
        };

        if (memoryAlarm.value) {
            options.value.annotations.yaxis.push(alarm(memoryAlarm.value));
        }

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [
            {
                name: 'Memory Max',
                data: memoryMaxSeries.value
            },
            {
                name: 'Memory Avg',
                data: memoryAvgSeries.value
            },
            {
                name: 'Memory Min',
                data: memoryMinSeries.value
            }
        ]);

        onMounted(() => {
            reloadInterval = setInterval(() => router.reload(), reloadTimer);

            if (!chartEl.value) {
                return;
            }

            nextTick(() => {
                chartEl.value.toggleSeries('Memory Max');
                chartEl.value.toggleSeries('Memory Min');
            });
        });

        onUnmounted(() => {
            clearInterval(reloadInterval);
        });

        return {
            subtitle,
            chartEl,
            options,
            links,
            series
        };
    }
};
</script>
