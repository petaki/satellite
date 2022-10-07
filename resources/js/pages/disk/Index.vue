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
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? `/disk?probe=${probe}&path=${diskPath}`
                                  : `/disk?probe=${probe}&path=${diskPath}&type=${type.value}`"
                              class="hover:text-cyan-500 sm:ml-3"
                              :class="{'text-cyan-500': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
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

<script>
import {
    CircleStackIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    toRefs,
    computed,
    onMounted,
    onUnmounted, nextTick
} from 'vue';

import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';
import useAnnotation from '../../common/useAnnotation';

export default {
    components: {
        CircleStackIcon,
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
    },

    setup(props) {
        const {
            diskPath,
            diskMinSeries,
            diskMaxSeries,
            diskAvgSeries,
            diskAlarm
        } = toRefs(props);

        const { alarm, max } = useAnnotation();
        const subtitle = ref(`Disk - ${diskPath.value}`);
        const chartEl = ref();
        const options = ref({});
        const reloadTimer = 60000;
        let reloadInterval;

        options.value.annotations = {
            yaxis: [
                max(Math.max(...diskMaxSeries.value.map(value => value.y)))
            ]
        };

        if (diskAlarm.value) {
            options.value.annotations.yaxis.push(alarm(diskAlarm.value));
        }

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [
            {
                name: `Disk Max - ${diskPath.value}`,
                data: diskMaxSeries.value
            },
            {
                name: `Disk Avg - ${diskPath.value}`,
                data: diskAvgSeries.value
            },
            {
                name: `Disk Min - ${diskPath.value}`,
                data: diskMinSeries.value
            }
        ]);

        onMounted(() => {
            reloadInterval = setInterval(() => Inertia.reload(), reloadTimer);

            nextTick(() => {
                chartEl.value.toggleSeries(`Disk Max - ${diskPath.value}`);
                chartEl.value.toggleSeries(`Disk Min - ${diskPath.value}`);
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
