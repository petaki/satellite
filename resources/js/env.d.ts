/// <reference types="vite/client" />

import type { Head, Link } from '@inertiajs/vue3';
import type ApexCharts from 'apexcharts';
import type { ApexOptions } from 'apexcharts';
import type {
    ComponentOptionsMixin, ComputedOptions, DefineComponent, MethodOptions
} from 'vue';
import type {
    VueApexChartsComponentMethods, VueApexChartsComponentProps
} from 'vue3-apexcharts';
import type AppTitle from './base/AppTitle.vue';
import type { SeriesType } from './types';

declare global {
    interface Window {
        isDark: () => boolean
        apexTheme: () => ApexOptions
        createApex: () => void
        Apex?: ApexOptions
    }
}

type Apexchart = DefineComponent<
    VueApexChartsComponentProps,
    VueApexChartsComponentMethods,
    {},
    ComputedOptions,
    MethodOptions,
    ComponentOptionsMixin,
    ComponentOptionsMixin,
    { mounted: (chart: ApexCharts, options: ApexOptions) => true }
>;

declare module 'vue' {
    interface HTMLAttributes {
        'scroll-region'?: boolean | string
    }

    interface GlobalComponents {
        Apexchart: Apexchart
        AppTitle: typeof AppTitle
        InertiaHead: typeof Head
        InertiaLink: typeof Link
    }
}

declare module '@inertiajs/core' {
    export interface InertiaConfig {
        sharedPageProps: {
            title: string
            suffix: string
            seriesButtons: string[]
            seriesTypes: SeriesType[]
            seriesType?: string
            probe?: string
            logPaths?: string[]
            isProbeActive?: boolean
            isCpuActive?: boolean
            isMemoryActive?: boolean
            isLoadActive?: boolean
            isDiskActive?: boolean
            isLogActive?: boolean
        }
    }
}

declare module '*.vue';

export {};
