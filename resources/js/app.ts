import {
    createInertiaApp, Head, Link
} from '@inertiajs/vue3';
import { defineAsyncComponent } from 'vue';
import AppTitle from './base/AppTitle.vue';

import '../css/app.css';

window.isDark = () => localStorage.theme === 'dark'
    || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);

window.apexTheme = () => (window.isDark()
    ? {
        chart: {
            foreColor: '#cad5e2'
        },
        colors: [
            '#fcc800',
            '#51a2ff',
            '#00d492',
            '#ff6467',
            '#ffdf20',
            '#5ee9b5'
        ],
        grid: {
            borderColor: '#45556c'
        }
    }
    : {
        chart: {
            foreColor: '#373d3f'
        },
        colors: [
            '#f6d757',
            '#2563EB',
            '#86EFAC',
            '#fecaca',
            '#fef08a',
            '#bbf7d0'
        ],
        grid: {
            borderColor: '#e0e0e0'
        }
    });

window.createApex = () => {
    const theme = window.apexTheme();

    window.Apex = Object.assign(window.Apex ?? {}, {
        chart: {
            animations: {
                enabled: false
            },
            ...theme.chart
        },
        markers: {
            size: 0
        },
        stroke: {
            width: 3
        },
        colors: theme.colors,
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
        grid: theme.grid
    });
};

window.createApex();

createInertiaApp({
    pages: './pages',
    withApp(app) {
        app
            .component('Apexchart', defineAsyncComponent(() => import('vue3-apexcharts')))
            .component('InertiaHead', Head)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle);
    }
});
