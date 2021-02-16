<template>
    <div class="memory__index layout__index">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                {{ $metaInfo.title }}
            </li>
        </ol>
        <div class="card">
            <div class="card-header">
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('grid-3x2')" />
                </svg>
                <span class="mr-auto">
                    {{ $metaInfo.title }}
                </span>
                <inertia-link v-for="(type, index) in seriesTypes"
                              :key="type.value"
                              :href="index === 0
                                  ? '/memory'
                                  : `/memory?type=${type.value}`"
                              class="ml-3"
                              :class="{'text-white': seriesType === type.value}">
                    {{ type.name }}
                </inertia-link>
            </div>
            <div class="chart_body card-body">
                <apexchart type="line"
                           :series="series"
                           height="100%"
                           :options="options" />
            </div>
        </div>
    </div>
</template>

<script>
import { Inertia } from '@inertiajs/inertia';
import Layout from '../../common/Layout.vue';

export default {
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

    layout: Layout,

    metaInfo() {
        return {
            title: 'Memory'
        };
    },

    data() {
        return {
            reloadInterval: undefined,
            reloadTimer: 60000,
            options: {}
        };
    },

    computed: {
        series() {
            return [{
                name: 'Memory',
                data: this.memorySeries
            }];
        }
    },

    mounted() {
        this.reloadInterval = setInterval(() => Inertia.reload(), this.reloadTimer);
    },

    beforeDestroy() {
        clearInterval(this.reloadInterval);
    }
};
</script>
