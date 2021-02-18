import Vue from 'vue';
import VueMeta from 'vue-meta';
import VueApexCharts from 'vue-apexcharts';
import { App, plugin } from '@inertiajs/inertia-vue';
import { InertiaProgress } from '@inertiajs/progress';

try {
    window.Popper = require('popper.js').default;
    window.$ = window.jQuery = require('jquery');

    require('bootstrap');
} catch (e) {}

InertiaProgress.init();

Vue.config.productionTip = false;
Vue.use(plugin);
Vue.use(VueMeta);
Vue.use(VueApexCharts);

window.Apex = {
    theme: {
        mode: 'dark'
    },

    chart: {
        background: 'transparent',
        foreColor: '#b4c9de'
    },

    colors: [
        '#00bcd4'
    ],

    grid: {
        borderColor: '#74818f'
    },

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

Vue.component('apexchart', VueApexCharts);

Vue.mixin({
    methods: {
        icon(name) {
            return `${this.$page.props.icons}#${name}`;
        }
    }
});

const el = document.getElementById('app');
const initialPage = JSON.parse(el.dataset.page);

new Vue({
    metaInfo: {
        titleTemplate: title => (title ? `${title} - ${initialPage.props.title}` : initialPage.props.title)
    },

    render: h => h(App, {
        props: {
            initialPage,
            resolveComponent: name => import(`./pages/${name}.vue`).then(module => module.default)
        }
    })
}).$mount(el);
