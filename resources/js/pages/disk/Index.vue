<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <circle-stack-icon class="h-6 w-6 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
                <span class="text-slate-600 sm:mx-auto">
                    Point Interval:
                    <span class="text-cyan-500">
                        {{ duration(chunkSize) }}
                    </span>
                </span>
                <select v-model="selectedType" class="form-select">
                    <option v-for="type in seriesTypes"
                            :key="type.value"
                            :value="type.value">
                        {{ type.name }}
                    </option>
                </select>
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

<script setup>
import {
    CircleStackIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    onMounted,
    onUnmounted,
    nextTick,
    defineProps,
    defineOptions, watch
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import useAnnotation from '../../base/useAnnotation';

const {
    probe,
    seriesType,
    seriesTypes,
    diskPath,
    diskMinSeries,
    diskMaxSeries,
    diskAvgSeries,
    diskAlarm
} = defineProps({
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

    diskPath: {
        type: String,
        required: true
    },

    diskMinSeries: {
        type: Array,
        default: () => []
    },

    diskMaxSeries: {
        type: Array,
        default: () => []
    },

    diskAvgSeries: {
        type: Array,
        default: () => []
    },

    diskAlarm: {
        type: Number,
        default: 0
    }
});

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const subtitle = ref(`Disk - ${diskPath}`);
const selectedType = ref(seriesType);
const chartEl = ref();
const reloadTimer = 60000;
let reloadInterval;

const options = ref({
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

options.value.annotations = {
    yaxis: [
        max(Math.max(...(diskMaxSeries ?? []).map(value => value.y)))
    ]
};

if (diskAlarm) {
    options.value.annotations.yaxis.push(alarm(diskAlarm));
}

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

onMounted(() => {
    reloadInterval = setInterval(() => router.reload(), reloadTimer);

    if (!chartEl.value) {
        return;
    }

    nextTick(() => {
        chartEl.value.toggleSeries(`Disk Max - ${diskPath}`);
        chartEl.value.toggleSeries(`Disk Min - ${diskPath}`);
    });
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});

watch(selectedType, () => {
    const index = seriesTypes.map(type => type.value).indexOf(selectedType.value);

    const href = index === 0
        ? `/disk?probe=${probe}&path=${diskPath}`
        : `/disk?probe=${probe}&path=${diskPath}&type=${selectedType.value}`;

    router.visit(href);
});
</script>
