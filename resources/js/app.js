import { DateTime, Duration } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, Head, Link } from '@inertiajs/vue3';
import VueApexCharts from 'vue3-apexcharts';
import AppTitle from './common/AppTitle.vue';

window.Apex = {
    chart: {
        animations: {
            enabled: false
        }
    },

    markers: {
        size: 0
    },

    stroke: {
        width: 3
    },

    colors: [
        '#f6d757',
        '#2563EB',
        '#86EFAC'
    ],

    tooltip: {
        x: {
            format: 'MMM. dd. HH:mm'
        }
    },

    xaxis: {
        type: 'datetime',
        labels: {
            datetimeUTC: false
        }
    },

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
};

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .use(VueApexCharts)
            .component('InertiaHead', Head)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle)
            .mixin({
                methods: {
                    date(value) {
                        const date = DateTime.fromISO(value);

                        if (!date.isValid) {
                            return value;
                        }

                        return date.toLocaleString(DateTime.DATETIME_MED);
                    },

                    duration(minutes) {
                        return Duration.fromObject({ minutes }).toFormat('hh:mm:ss');
                    }
                }
            })
            .mount(el);
    }
});
