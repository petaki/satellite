import {
    createInertiaApp, Head, Link
} from '@inertiajs/vue3';
import VueApexCharts from 'vue3-apexcharts';
import type { ApexConfig } from './types';
import AppTitle from './base/AppTitle.vue';

import '../css/app.css';

window.isDark = () => localStorage.theme === 'dark'
    || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);

window.createApex = () => {
    const apex: ApexConfig = {
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
    pages: './pages',
    withApp(app) {
        app
            .use(VueApexCharts)
            .component('InertiaHead', Head)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle);
    }
});
