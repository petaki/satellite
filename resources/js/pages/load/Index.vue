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
                <apexchart v-if="load1Series"
                           ref="chartEl"
                           type="line"
                           :series="series"
                           height="100%"
                           :options="options" />
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
    watch, nextTick
} from 'vue';

import Breadcrumb from '../../base/Breadcrumb.vue';
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
    load1Series = [],
    load5Series = [],
    load15Series = [],
    loadAlarm = 0
} = defineProps<{
    probe: string
    chunkSize: number
    load1Series?: SeriesDataPoint[]
    load5Series?: SeriesDataPoint[]
    load15Series?: SeriesDataPoint[]
    loadAlarm?: number
}>();

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const { isLive } = useLiveReload();
const subtitle = ref('Load');
const chartEl = ref();

const series = computed(() => [
    {
        name: 'Load 1',
        data: load1Series
    },
    {
        name: 'Load 5',
        data: load5Series
    },
    {
        name: 'Load 15',
        data: load15Series
    }
]);

const options = ref<ApexConfig>({
    dataLabels: {
        enabled: false
    },
    yaxis: {
        min: 0,
        labels: {
            formatter(val: number) {
                if (val) {
                    return val.toFixed(2);
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
    ? `/load?probe=${probe}`
    : `/load?probe=${probe}&type=${selectedType}`);

const onSetTheme = () => {
    chartEl.value.refresh();
};

const refreshSeries = () => {
    if (!chartEl.value) {
        return;
    }

    if (loadAlarm) {
        options.value.annotations = {
            yaxis: [
                alarm(loadAlarm, '')
            ]
        };
    }

    nextTick(() => {
        chartEl.value.refresh();
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
