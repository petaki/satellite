<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <circle-stack-icon class="h-6 w-6 sm:mr-2" />
                <span class="flex-1 sm:mr-auto">
                    {{ subtitle }}
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
                <apexchart v-if="diskSeries"
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
    onUnmounted
} from 'vue';

import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

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

        diskSeries: {
            type: Array,
            default: () => []
        },

        diskAlarm: {
            type: Number,
            default: 0
        }
    },

    setup(props) {
        const { diskPath, diskSeries, diskAlarm } = toRefs(props);
        const subtitle = ref(`Disk - ${diskPath.value}`);
        const reloadTimer = 60000;

        let reloadInterval;

        const options = ref(diskAlarm.value
            ? {
                annotations: {
                    yaxis: [
                        {
                            y: diskAlarm.value,
                            borderColor: '#ef4444',
                            label: {
                                borderColor: '#ef4444',
                                style: {
                                    color: '#fff',
                                    background: '#ef4444'
                                },
                                text: `Alarm: ${diskAlarm.value}%`
                            }
                        }
                    ]
                }
            }
            : {});

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [{
            name: `Disk - ${diskPath.value}`,
            data: diskSeries.value
        }]);

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
