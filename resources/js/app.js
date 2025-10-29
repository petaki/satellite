import { DateTime, Duration } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, Head, Link } from '@inertiajs/vue3';
import VueApexCharts from 'vue3-apexcharts';
import AppTitle from './base/AppTitle.vue';

import '../css/app.css';

window.isDark = () => localStorage.theme === 'dark'
    || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);

window.createApex = () => {
    const apex = {
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
            '#86EFAC',
            '#fecaca',
            '#fef08a',
            '#bbf7d0'
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
        }
    };

    if (window.isDark()) {
        apex.chart.foreColor = '#cad5e2';

        apex.colors = [
            '#fcc800',
            '#51a2ff',
            '#00d492',
            '#ff6467',
            '#ffdf20',
            '#5ee9b5'
        ];

        apex.grid = {
            borderColor: '#45556c'
        };
    }

    window.Apex = apex;
};

window.createApex();

createInertiaApp({
    resolve: name => {
        const pages = import.meta.glob('./pages/**/*.vue', { eager: true });

        return pages[`./pages/${name}.vue`];
    },
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
