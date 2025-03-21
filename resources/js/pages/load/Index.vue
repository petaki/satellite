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

<script setup>
import {
    DocumentDuplicateIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    defineProps,
    defineOptions, watch, nextTick
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import useAnnotation from '../../base/useAnnotation';
import SeriesSelector from '../../base/SeriesSelector.vue';

const {
    probe,
    load1Series,
    load5Series,
    load15Series,
    loadAlarm
} = defineProps({
    probe: {
        type: String,
        required: true
    },

    chunkSize: {
        type: Number,
        required: true
    },

    load1Series: {
        type: Array,
        default: () => []
    },

    load5Series: {
        type: Array,
        default: () => []
    },

    load15Series: {
        type: Array,
        default: () => []
    },

    loadAlarm: {
        type: Number,
        default: 0
    }
});

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const subtitle = ref('Load');
const chartEl = ref();
const reloadTimer = 60000;
let reloadInterval;

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

const options = ref({
    dataLabels: {
        enabled: false
    },
    yaxis: {
        min: 0,
        labels: {
            formatter(val) {
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

const seriesHref = (isDefault, selectedType) => (isDefault
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
    reloadInterval = setInterval(() => { router.reload(); }, reloadTimer);

    refreshSeries();

    document.addEventListener('set-theme', onSetTheme);
});

onUnmounted(() => {
    clearInterval(reloadInterval);

    document.removeEventListener('set-theme', onSetTheme);
});
</script>
