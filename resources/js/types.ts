import type { Ref } from 'vue';

export interface BreadcrumbLink {
    name: string | Ref<string>
    href?: string
}

export interface SeriesDataPoint {
    x: number
    y: number
    name?: string
}

export interface SeriesType {
    name: string
    value: string
}

export interface YAxisAnnotation {
    y: number
    borderColor: string
    label: {
        borderColor: string
        style: {
            color: string
            fontSize: string
            fontWeight: number
            background: string
        }
        text: string
        offsetY: number
    }
}

export interface ProbeSummary {
    name: string
    cpu: number
    memory: number
    load1: number
    load5: number
    load15: number
    isActive: boolean
    cpuAlarm: number
    memAlarm: number
    loadAlarm: number
}

export interface ApexConfig {
    chart?: {
        animations?: {
            enabled: boolean
        }
        foreColor?: string
        stacked?: boolean
    }
    markers?: {
        size: number
    }
    stroke?: {
        width: number
    }
    colors?: string[]
    tooltip?: {
        x?: {
            format: string
        }
        y?: {
            formatter: (
                value: number,
                opts: { seriesIndex: number, dataPointIndex: number }
            ) => string
        }
        marker?: {
            show: boolean
        }
    }
    xaxis?: {
        type: string
        labels?: {
            datetimeUTC: boolean
        }
    }
    yaxis?: {
        min?: number
        max?: number
        forceNiceScale?: boolean
        labels?: {
            formatter: (val: number) => string
        }
    }
    grid?: {
        borderColor: string
    }
    dataLabels?: {
        enabled: boolean
    }
    annotations?: {
        yaxis: YAxisAnnotation[]
    }
}
