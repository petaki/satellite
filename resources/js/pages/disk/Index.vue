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
                <circle-stack-icon class="h-6 w-6 shrink-0" />
                <!-- eslint-disable max-len -->
                <select class="form-select flex-1 sm:mx-2"
                        :value="diskPath || ''"
                        @change="onPathChange">
                    <option v-for="path in diskPaths"
                            :key="path"
                            :value="path">
                        {{ path }}
                    </option>
                </select>
                <!-- eslint-enable max-len -->
                <PointInterval v-if="diskPath" :chunk-size="chunkSize" />
                <SeriesSelector v-if="diskPath" :href="seriesHref" />
            </card-title>
            <div v-if="diskPath" class="chart">
                <apexchart v-if="diskMinSeries"
                           ref="chartEl"
                           type="line"
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
    CircleStackIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    nextTick,
    watch
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import ChartEmpty from '../../base/ChartEmpty.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import LiveIndicator from '../../base/LiveIndicator.vue';
import RefreshButton from '../../base/RefreshButton.vue';
import useAnnotation from '../../use/useAnnotation';
import useLiveReload from '../../use/useLiveReload';
import usePaths from '../../use/usePaths';
import PointInterval from '../../base/PointInterval.vue';
import SeriesSelector from '../../base/SeriesSelector.vue';
import type { SeriesDataPoint, ApexConfig } from '../../types';

const {
    probe,
    chunkSize,
    diskPath = '',
    diskPaths = [],
    diskMinSeries = [],
    diskMaxSeries = [],
    diskAvgSeries = [],
    diskAlarm = 0
} = defineProps<{
    probe: string
    chunkSize: number
    diskPath?: string
    diskPaths?: string[]
    diskMinSeries?: SeriesDataPoint[]
    diskMaxSeries?: SeriesDataPoint[]
    diskAvgSeries?: SeriesDataPoint[]
    diskAlarm?: number
}>();

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const { isLive } = useLiveReload();
const { diskPath: diskPathFn } = usePaths();
const subtitle = ref('Disk');

const onPathChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;

    router.visit(diskPathFn(probe, target.value));
};
const chartEl = ref();

const options = ref<ApexConfig>({
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

const series = computed(() => [
    {
        name: `Disk Max - ${diskPath}`,
        data: diskMaxSeries
    },
    {
        name: `Disk Avg - ${diskPath}`,
        data: diskAvgSeries
    },
    {
        name: `Disk Min - ${diskPath}`,
        data: diskMinSeries
    }
]);

const seriesHref = (isDefault: boolean, selectedType: string | undefined) => (isDefault
    ? diskPathFn(probe, diskPath)
    : diskPathFn(probe, diskPath, selectedType));

const onSetTheme = () => {
    chartEl.value.refresh();
};

const refreshSeries = async () => {
    if (!chartEl.value) {
        return;
    }

    options.value.annotations = {
        yaxis: [
            max(Math.max(...(diskMaxSeries ?? []).map(value => value.y)))
        ]
    };

    if (diskAlarm) {
        options.value.annotations.yaxis.push(alarm(diskAlarm));
    }

    await nextTick(() => {
        chartEl.value.refresh();
    });

    await nextTick(() => {
        chartEl.value.toggleSeries(`Disk Max - ${diskPath}`);
        chartEl.value.toggleSeries(`Disk Min - ${diskPath}`);
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
