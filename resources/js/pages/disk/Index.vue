<template>
    <inertia-head :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <database-icon class="h-6 w-6 mr-2" />
                <span class="flex-1 mr-auto">
                    {{ subtitle }}
                </span>
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? `/disk?path=${diskPath}`
                                  : `/disk?path=${diskPath}&type=${type.value}`"
                              class="hover:text-cyan-500 ml-3"
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
    DatabaseIcon
} from '@heroicons/vue/outline';

import {
    ref, toRefs, computed, onMounted, onUnmounted
} from 'vue';
import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        DatabaseIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
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
        }
    },

    setup(props) {
        const { diskPath, diskSeries } = toRefs(props);
        const subtitle = ref(`Disk - ${diskPath.value}`);
        const reloadInterval = ref();
        const reloadTimer = ref(60000);
        const options = ref({});

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [{
            name: `Disk - ${diskPath.value}`,
            data: diskSeries.value
        }]);

        onMounted(() => {
            reloadInterval.value = setInterval(() => Inertia.reload(), reloadTimer.value);
        });

        onUnmounted(() => {
            clearInterval(reloadInterval.value);
        });

        return {
            subtitle,
            reloadInterval,
            reloadTimer,
            options,
            links,
            series
        };
    }
};
</script>
