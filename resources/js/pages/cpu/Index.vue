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
                <cpu-chip-icon class="h-6 w-6 shrink-0 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
                <PointInterval :chunk-size="chunkSize" />
                <SeriesSelector :href="seriesHref" />
            </card-title>
            <div class="chart">
                <apexchart v-if="cpuMinSeries.length"
                           ref="chartEl"
                           type="bar"
                           :series="series"
                           height="100%"
                           :options="options"
                           @mounted="onChartMounted" />
                <chart-empty v-else />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    CpuChipIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    nextTick,
    onMounted,
    onUnmounted,
    watch
} from 'vue';

import type { ApexOptions } from 'apexcharts';

import Breadcrumb from '../../base/Breadcrumb.vue';
import ChartEmpty from '../../base/ChartEmpty.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import LiveIndicator from '../../base/LiveIndicator.vue';
import RefreshButton from '../../base/RefreshButton.vue';
import useAnnotation from '../../use/useAnnotation';
import escapeHtml from '../../use/useHtml';
import useLiveReload from '../../use/useLiveReload';
import usePaths from '../../use/usePaths';
import PointInterval from '../../base/PointInterval.vue';
import SeriesSelector from '../../base/SeriesSelector.vue';
import type { SeriesDataPoint } from '../../types';

const {
    probe,
    chunkSize,
    cpuMinSeries = [],
    cpuMaxSeries = [],
    cpuAvgSeries = [],
    process1Series = [],
    process2Series = [],
    process3Series = [],
    cpuAlarm = 0
} = defineProps<{
    probe: string
    chunkSize: number
    cpuMinSeries?: SeriesDataPoint[]
    cpuMaxSeries?: SeriesDataPoint[]
    cpuAvgSeries?: SeriesDataPoint[]
    process1Series?: SeriesDataPoint[]
    process2Series?: SeriesDataPoint[]
    process3Series?: SeriesDataPoint[]
    cpuAlarm?: number
}>();

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const { isLive } = useLiveReload();
const { cpuPath } = usePaths();
const subtitle = ref('CPU');
const chartEl = ref();

const series = computed(() => [
    {
        name: 'CPU Max',
        type: 'line',
        data: cpuMaxSeries
    },
    {
        name: 'CPU Avg',
        type: 'line',
        data: cpuAvgSeries
    },
    {
        name: 'CPU Min',
        type: 'line',
        data: cpuMinSeries
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

const options = ref<ApexOptions>({
    plotOptions: {
        bar: {
            columnWidth: '70%'
        }
    },
    chart: {
        stacked: true
    },
    dataLabels: {
        enabled: false
    },
    tooltip: {
        y: {
            formatter(value, opts) {
                const percent = value === null ? '' : `${value.toFixed(2)}%`;
                const name = opts && opts.seriesIndex > 2
                    ? series.value[opts.seriesIndex].data[opts.dataPointIndex].name
                    : undefined;

                if (name !== undefined) {
                    return `${escapeHtml(name)}: ${percent}`;
                }

                return percent;
            }
        },
        marker: {
            show: false
        }
    },
    yaxis: {
        forceNiceScale: true,
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

const seriesHref = (selectedType: string) => cpuPath(probe, selectedType);

const onSetTheme = () => {
    chartEl.value?.updateOptions(window.apexTheme());
};

const refreshSeries = async () => {
    await nextTick();

    if (!chartEl.value?.chart) {
        return;
    }

    const yaxis = [
        max(Math.max(...(cpuMaxSeries ?? []).map(value => value.y)))
    ];

    if (cpuAlarm) {
        yaxis.push(alarm(cpuAlarm));
    }

    options.value.annotations = { yaxis };

    chartEl.value.updateOptions({ annotations: options.value.annotations });
};

const onChartMounted = async () => {
    await refreshSeries();

    chartEl.value?.hideSeries('CPU Max');
    chartEl.value?.hideSeries('CPU Min');
};

watch(series, refreshSeries, { flush: 'post' });

onMounted(() => {
    document.addEventListener('set-theme', onSetTheme);
});

onUnmounted(() => {
    document.removeEventListener('set-theme', onSetTheme);
});
</script>
