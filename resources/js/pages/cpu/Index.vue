<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <cpu-chip-icon class="h-6 w-6 sm:mr-2" />
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
                                  ? `/cpu?probe=${probe}`
                                  : `/cpu?probe=${probe}&type=${type.value}`"
                              class="hover:text-cyan-500 sm:ml-3"
                              :class="{'text-cyan-500': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
            </card-title>
            <div class="chart">
                <apexchart v-if="cpuMinSeries"
                           ref="chartEl"
                           type="bar"
                           :series="series"
                           height="100%"
                           :options="options" />
            </div>
        </div>
    </div>
</template>

<script>
import {
    CpuChipIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    toRefs,
    computed,
    nextTick,
    onMounted,
    onUnmounted
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';
import useAnnotation from '../../common/useAnnotation';

export default {
    components: {
        CpuChipIcon,
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

        cpuMinSeries: {
            type: Array,
            default: () => []
        },

        cpuMaxSeries: {
            type: Array,
            default: () => []
        },

        cpuAvgSeries: {
            type: Array,
            default: () => []
        },

        cpuAlarm: {
            type: Number,
            default: 0
        }
    },

    setup(props) {
        const {
            cpuMinSeries,
            cpuMaxSeries,
            cpuAvgSeries,
            cpuAlarm
        } = toRefs(props);

        const { alarm, max } = useAnnotation();
        const subtitle = ref('CPU');
        const chartEl = ref();
        const reloadTimer = 60000;
        let reloadInterval;

        const series = computed(() => [
            {
                name: 'CPU Max',
                type: 'line',
                data: cpuMaxSeries.value
            },
            {
                name: 'CPU Avg',
                type: 'line',
                data: cpuAvgSeries.value
            },
            {
                name: 'CPU Min',
                type: 'line',
                data: cpuMinSeries.value
            },
            {
                name: 'Process #1',
                type: 'column',
                data: []
            },
            {
                name: 'Process #2',
                type: 'column',
                data: []
            },
            {
                name: 'Process #3',
                type: 'column',
                data: []
            }
        ]);

        cpuAvgSeries.value.forEach(value => {
            if (!value.p) {
                return;
            }

            value.p.forEach((process, index) => {
                if (index > 2) {
                    return;
                }

                series.value[index + 3].data.push({
                    x: value.x, ...process
                });
            });
        });

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
                }
            }
        });

        options.value.annotations = {
            yaxis: [
                max(Math.max(...(cpuMaxSeries.value ?? []).map(value => value.y)))
            ]
        };

        if (cpuAlarm.value) {
            options.value.annotations.yaxis.push(alarm(cpuAlarm.value));
        }

        const links = ref([
            { name: subtitle }
        ]);

        onMounted(() => {
            reloadInterval = setInterval(() => router.reload(), reloadTimer);

            if (!chartEl.value) {
                return;
            }

            nextTick(() => {
                chartEl.value.toggleSeries('CPU Max');
                chartEl.value.toggleSeries('CPU Min');
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
