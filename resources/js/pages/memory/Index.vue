<template>
    <inertia-head :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <duplicate-icon class="h-6 w-6 mr-2" />
                <span class="flex-1 mr-auto">
                    {{ subtitle }}
                </span>
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? '/memory'
                                  : `/memory?type=${type.value}`"
                              class="hover:text-cyan-500 ml-3"
                              :class="{'text-cyan-500': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
            </card-title>
            <div class="chart">
                <apexchart v-if="memorySeries"
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
    DuplicateIcon
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
        DuplicateIcon,
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

        memorySeries: {
            type: Array,
            default: () => []
        }
    },

    setup(props) {
        const { memorySeries } = toRefs(props);
        const subtitle = ref('Memory');
        const reloadInterval = ref();
        const reloadTimer = ref(60000);
        const options = ref({});

        const links = ref([
            { name: subtitle }
        ]);

        const series = computed(() => [{
            name: 'Memory',
            data: memorySeries.value
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
