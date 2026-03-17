<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8 dark:bg-slate-700">
            <card-title>
                <circle-stack-icon class="h-6 w-6 sm:mr-2" />
                <span>
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
                <apexchart v-if="diskMinSeries"
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
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import useAnnotation from '../../use/useAnnotation';
import useDate from '../../use/useDate';
import SeriesSelector from '../../base/SeriesSelector.vue';
import type { SeriesDataPoint, ApexConfig } from '../../types';

const {
    probe,
    chunkSize,
    diskPath,
    diskMinSeries = [],
    diskMaxSeries = [],
    diskAvgSeries = [],
    diskAlarm = 0
} = defineProps<{
    probe: string
    chunkSize: number
    diskPath: string
    diskMinSeries?: SeriesDataPoint[]
    diskMaxSeries?: SeriesDataPoint[]
    diskAvgSeries?: SeriesDataPoint[]
    diskAlarm?: number
}>();

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const { duration } = useDate();
const subtitle = ref(`Disk - ${diskPath}`);
const chartEl = ref();
const reloadTimer = 60000;
let reloadInterval: ReturnType<typeof setInterval>;

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
    ? `/disk?probe=${probe}&path=${diskPath}`
    : `/disk?probe=${probe}&path=${diskPath}&type=${selectedType}`);

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
    reloadInterval = setInterval(() => { router.reload(); }, reloadTimer);

    refreshSeries();

    document.addEventListener('set-theme', onSetTheme);
});

onUnmounted(() => {
    clearInterval(reloadInterval);

    document.removeEventListener('set-theme', onSetTheme);
});
</script>
