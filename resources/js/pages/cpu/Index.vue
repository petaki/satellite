<template>
    <div class="cpu__index layout__index">
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
                    <use :xlink:href="icon('cpu')" />
                </svg>
                <span>
                    {{ $metaInfo.title }}
                </span>
            </div>
            <div class="chart_body card-body">
                <div ref="chart"></div>
            </div>
        </div>
    </div>
</template>

<script>
import ApexCharts from 'apexcharts';
import Layout from '../../common/Layout.vue';

export default {
    props: {
        cpuSeries: {
            type: Array,
            default: () => []
        }
    },

    layout: Layout,

    metaInfo() {
        return {
            title: 'CPU'
        };
    },

    mounted() {
        const options = {
            theme: {
                mode: 'dark'
            },

            chart: {
                type: 'line',
                height: '100%',
                background: 'transparent'
            },

            colors: [
                '#00bcd4'
            ],

            grid: {
                borderColor: '#74818f'
            },

            series: [{
                name: this.$metaInfo.title,
                data: this.cpuSeries
            }],

            xaxis: {
                type: 'datetime'
            },

            yaxis: {
                labels: {
                    formatter(val, index) {
                        return `${val.toFixed(2)}%`;
                    }
                }
            }
        };

        const chart = new ApexCharts(this.$refs.chart, options);

        chart.render();
    }
};
</script>
