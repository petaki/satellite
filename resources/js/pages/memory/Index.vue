<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8 dark:bg-slate-700">
            <card-title>
                <document-duplicate-icon class="h-6 w-6 sm:mr-2 dark:text-slate-300" />
                <span class="dark:text-slate-300">
                    {{ subtitle }}
                </span>
                <span class="text-slate-600 sm:mx-auto dark:text-slate-500">
                    Point Interval:
                    <span class="text-cyan-500">
                        {{ duration(chunkSize) }}
                    </span>
                </span>
                <SeriesSelector :href="seriesHref" />
            </card-title>
            <div class="chart">
                <apexchart v-if="memoryMinSeries"
                           ref="chartEl"
                           type="bar"
                           :series="series"
                           height="100%"
                           :options="options" />
            </div>
        </div>
    </div>
</template>

<script setup>
import {
    DocumentDuplicateIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    nextTick,
    defineProps,
    defineOptions,
    watch
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import useAnnotation from '../../base/useAnnotation';
import SeriesSelector from '../../base/SeriesSelector.vue';

const {
    probe,
    memoryMinSeries,
    memoryMaxSeries,
    memoryAvgSeries,
    process1Series,
    process2Series,
    process3Series,
    memoryAlarm
} = defineProps({
    probe: {
        type: String,
        required: true
    },

    chunkSize: {
        type: Number,
        required: true
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

    process1Series: {
        type: Array,
        default: () => []
    },

    process2Series: {
        type: Array,
        default: () => []
    },

    process3Series: {
        type: Array,
        default: () => []
    },

    memoryAlarm: {
        type: Number,
        default: 0
    }
});

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const subtitle = ref('Memory');
const chartEl = ref();
const reloadTimer = 60000;
let reloadInterval;

const series = computed(() => [
    {
        name: 'Memory Max',
        type: 'line',
        data: memoryMaxSeries
    },
    {
        name: 'Memory Avg',
        type: 'line',
        data: memoryAvgSeries
    },
    {
        name: 'Memory Min',
        type: 'line',
        data: memoryMinSeries
    },
    {
        name: 'Process #1',
        type: 'column',
        data: process1Series
    },
    {
        name: 'Process #2',
        type: 'column',
        data: process2Series
    },
    {
        name: 'Process #3',
        type: 'column',
        data: process3Series
    }
]);

const options = ref({
    chart: {
        stacked: true
    },
    dataLabels: {
        enabled: false
    },
    tooltip: {
        y: {
            formatter(value, { seriesIndex, dataPointIndex }) {
                if (seriesIndex > 2) {
                    return `${series.value[seriesIndex].data[dataPointIndex].name}: ${value.toFixed(2)}%`;
                }

                return `${value.toFixed(2)}%`;
            }
        },
        marker: {
            show: false
        }
    },
    yaxis: {
        min: 0,
        max: 100,
        labels: {
            formatter(val) {
                if (val) {
                    return `${val.toFixed(2)}%`;
                }

                return '';
            }
        }
    }
});

const links = ref([
    { name: subtitle }
]);

const seriesHref = (isDefault, selectedType) => (isDefault
    ? `/memory?probe=${probe}`
    : `/memory?probe=${probe}&type=${selectedType}`);

const onSetTheme = () => {
    chartEl.value.refresh();
};

const refreshSeries = async () => {
    if (!chartEl.value) {
        return;
    }

    options.value.annotations = {
        yaxis: [
            max(Math.max(...(memoryMaxSeries ?? []).map(value => value.y)))
        ]
    };

    if (memoryAlarm) {
        options.value.annotations.yaxis.push(alarm(memoryAlarm));
    }

    await nextTick(() => {
        chartEl.value.refresh();
    });

    await nextTick(() => {
        chartEl.value.toggleSeries('Memory Max');
        chartEl.value.toggleSeries('Memory Min');
    });
};

watch(series, refreshSeries);

onMounted(() => {
    reloadInterval = setInterval(() => router.reload(), reloadTimer);

    refreshSeries();

    document.addEventListener('set-theme', onSetTheme);
});

onUnmounted(() => {
    clearInterval(reloadInterval);

    document.removeEventListener('set-theme', onSetTheme);
});
</script>
