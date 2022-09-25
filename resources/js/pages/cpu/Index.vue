<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <cpu-chip-icon class="h-6 w-6 sm:mr-2" />
                <span class="flex-1 sm:mr-auto">
                    {{ subtitle }}
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
    CpuChipIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    toRefs,
    computed,
    onMounted,
    onUnmounted
} from 'vue';

import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

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

        const subtitle = ref('CPU');
        const options = ref({});
        const reloadTimer = 60000;
        let reloadInterval;

        if (cpuAlarm.value) {
            options.value.annotations = {
                yaxis: [
                    {
                        y: cpuAlarm.value,
                        borderColor: '#ef4444',
                        label: {
                            borderColor: '#ef4444',
                            style: {
                                color: '#fff',
                                background: '#ef4444'
                            },
                            text: `Alarm: ${cpuAlarm.value}%`
                        }
                    }
                ]
            };
        }

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [
            {
                name: 'CPU Max',
                data: cpuMaxSeries.value
            },
            {
                name: 'CPU Avg',
                data: cpuAvgSeries.value
            },
            {
                name: 'CPU Min',
                data: cpuMinSeries.value
            }
        ]);

        onMounted(() => {
            reloadInterval = setInterval(() => Inertia.reload(), reloadTimer);
        });

        onUnmounted(() => {
            clearInterval(reloadInterval);
        });

        return {
            subtitle,
            options,
            links,
            series
        };
    }
};
</script>
