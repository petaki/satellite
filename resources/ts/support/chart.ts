import BaseChart, { ChartData, ChartOptions } from 'chart.js';
import _ from 'lodash';
import moment from 'moment';
import { ChartType } from '../store/types';
import { IDataset } from './types';

export default class Chart {
    protected isInitialized = false;
    protected context: HTMLCanvasElement;
    protected chart: BaseChart;

    constructor(context: HTMLCanvasElement) {
        this.context = context;
    }

    public init(chartType: ChartType, dataset: IDataset): void {
        const data = this.data(chartType, dataset);
        const options = this.options();

        this.chart = new BaseChart(this.context, {
            data,
            options,
            type: 'line',
        });

        this.isInitialized = true;
    }

    public update(chartType: ChartType, dataset: IDataset): void {
        if (!this.isInitialized) {
            return;
        }

        this.chart.data = this.data(chartType, dataset);
        this.chart.update();
    }

    protected data(chartType: ChartType, dataset: IDataset): ChartData {
        const backgroundColor = '#3490dc';
        const borderColor = '#3490dc';

        const labels =  _.map(
            _.keys(dataset), (key) => moment.unix(+key).format(this.format(chartType)),
        );

        const data = _.map(
            _.values(dataset), parseFloat,
        );

        return {
            labels,

            datasets: [{
                backgroundColor,
                borderColor,
                borderWidth: 1,
                data,
                fill: false,
                lineTension: 0,
                pointHoverRadius: 0,
                pointRadius: 0,
            }],
        };
    }

    protected format(chartType: ChartType): string {
        if (chartType === ChartType.Day) {
            return 'HH:mm';
        }

        if (chartType === ChartType.Week) {
            return 'ddd, HH:mm';
        }

        return 'MMM DD, HH:mm';
    }

    protected options(): ChartOptions {
        const color = '#2e3856';
        const fontColor = '#596b9f';

        return {
            responsive: true,

            legend: {
                display: false,
            },
            scales: {
                xAxes: [{
                    gridLines: {
                        color,
                    },
                    ticks: {
                        fontColor,
                    },
                }],
                yAxes: [{
                    gridLines: {
                        color,
                    },
                    ticks: {
                        beginAtZero: true,
                        fontColor,

                        callback(value) {
                            return `${value}%`;
                        },
                    },
                }],
            },
            tooltips: {
                intersect: false,

                callbacks: {
                    label(tooltipItems) {
                        if (_.isUndefined(tooltipItems.yLabel)) {
                            return '0%';
                        }

                        const value = parseFloat(tooltipItems.yLabel).toFixed(2);

                        return `${value}%`;
                    },
                },
            },
        };
    }
}
