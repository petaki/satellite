import { DateTime } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, InertiaLink } from '@inertiajs/inertia-vue3';
import { InertiaProgress } from '@inertiajs/progress';
import VueApexCharts from 'vue3-apexcharts';
import AppTitle from './common/AppTitle.vue';

window.Apex = {
    chart: {
        animations: {
            enabled: false
        }
    },

    colors: [
        'rgb(6, 182, 212)'
    ],

    xaxis: {
        type: 'datetime'
    },

    yaxis: {
        max: 100,
        labels: {
            formatter(val) {
                return `${val.toFixed(2)}%`;
            }
        }
    }
};

InertiaProgress.init();

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .use(VueApexCharts)
            .component('InertiaLink', InertiaLink)
            .component('AppTitle', AppTitle)
            .mixin({
                methods: {
                    date(value) {
                        const date = DateTime.fromISO(value);

                        if (!date.isValid) {
                            return value;
                        }

                        return date.toLocaleString(DateTime.DATETIME_MED);
                    }
                }
            })
            .mount(el);
    }
});
