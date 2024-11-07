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
                <select v-model="selectedType" class="form-select">
                    <option v-for="type in seriesTypes"
                            :key="type.value"
                            :value="type.value">
                        {{ type.name }}
                    </option>
                </select>
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
    load1Series,
    load5Series,
    load15Series
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
    }
});

defineOptions({
    layout: Layout
});

const { alarm, max } = useAnnotation();
const subtitle = ref('Load');
const selectedType = ref(seriesType);
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

onMounted(() => {
    reloadInterval = setInterval(() => router.reload(), reloadTimer);
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});

watch(selectedType, () => {
    const index = seriesTypes.map(type => type.value).indexOf(selectedType.value);

    const href = index === 0
        ? `/load?probe=${probe}`
        : `/load?probe=${probe}&type=${selectedType.value}`;

    router.visit(href);
});
</script>
