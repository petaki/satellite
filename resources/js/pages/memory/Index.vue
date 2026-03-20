<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links">
            <div class="flex items-center gap-2">
                <refresh-button />
                <live-indicator v-model="isLive" />
            </div>
        </breadcrumb>
        <div class="bg-white p-8 dark:bg-slate-700">
            <card-title>
                <document-duplicate-icon class="h-6 w-6 shrink-0 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
                <PointInterval :chunk-size="chunkSize" />
                <SeriesSelector :href="seriesHref" />
            </card-title>
            <div class="chart">
                <apexchart v-if="memoryMinSeries"
                           ref="chartEl"
                           type="bar"
                           :series="series"
                           height="100%"
                           :options="options" />
                <chart-empty v-else />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    DocumentDuplicateIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    nextTick,
    watch
} from 'vue';

import Breadcrumb from '../../base/Breadcrumb.vue';
import ChartEmpty from '../../base/ChartEmpty.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import LiveIndicator from '../../base/LiveIndicator.vue';
import RefreshButton from '../../base/RefreshButton.vue';
import useAnnotation from '../../use/useAnnotation';
import useLiveReload from '../../use/useLiveReload';
import PointInterval from '../../base/PointInterval.vue';
import SeriesSelector from '../../base/SeriesSelector.vue';
import type { SeriesDataPoint, ApexConfig } from '../../types';

const {
    probe,
    chunkSize,
    memoryMinSeries = [],
    memoryMaxSeries = [],
    memoryAvgSeries = [],
    process1Series = [],
    process2Series = [],
    process3Series = [],
    memoryAlarm = 0
} = defineProps<{
    probe: string
    chunkSize: number
    memoryMinSeries?: SeriesDataPoint[]
    memoryMaxSeries?: SeriesDataPoint[]
    memoryAvgSeries?: SeriesDataPoint[]
    process1Series?: SeriesDataPoint[]
    process2Series?: SeriesDataPoint[]
    process3Series?: SeriesDataPoint[]
    memoryAlarm?: number
}>();

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const { isLive } = useLiveReload();
const subtitle = ref('Memory');
const chartEl = ref();

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

const options = ref<ApexConfig>({
    chart: {
        stacked: true
    },
    dataLabels: {
        enabled: false
    },
    tooltip: {
        y: {
            formatter(value: number, { seriesIndex, dataPointIndex }) {
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
            formatter(val: number) {
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

const seriesHref = (isDefault: boolean, selectedType: string | undefined) => (isDefault
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
    refreshSeries();

    document.addEventListener('set-theme', onSetTheme);
});

onUnmounted(() => {
    document.removeEventListener('set-theme', onSetTheme);
});
</script>
