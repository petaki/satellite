<template>
    <div class="content">
        <div class="panel">
            <div class="flex items-center mb-6">
                <h4 class="text-indigo-lighter">
                    {{ selected.name }}
                </h4>
                <a v-for="(value, name, index) in types"
                   :key="value"
                   class="tab"
                   :class="{'ml-auto rounded-l': index === 0, 'rounded-r': index === typeSize - 1, active: chartType === value}"
                   href="#"
                   @click.prevent="editChartType(value)">
                    {{ name }}
                </a>
            </div>
            <canvas ref="chart"></canvas>
        </div>
    </div>
</template>
<script lang="ts">
    import _ from 'lodash';
    import BaseChart from '../support/chart';
    import Repository from '../support/repository';
    import { ChartType, IConnection, IDataset, ISelected, IFlash, FlashType } from '../store/types';
    import { Component, Vue, Watch } from 'vue-property-decorator';
    import { Action, State, Mutation } from 'vuex-class';

    @Component({
        name: 'Chart',
    })
    export default class Chart extends Vue {
        chart!: BaseChart;
        repository!: Repository;

        @State chartType!: ChartType;
        @State connection?: IConnection;
        @State selected!: ISelected;

        @Mutation editChartType!: (chartType: ChartType) => void;
        @Action flash!: (flash?: IFlash) => void;

        get types(): object {
            return _.pickBy(ChartType, (value) => _.isNumber(value));
        }

        get typeSize(): number {
            return _.size(this.types);
        }

        @Watch('chartType')
        @Watch('selected')
        onUpdateChart(): void {
            this.findDataset().then(dataset => {
                this.chart.update(this.chartType, dataset);
                this.flash();
            });
        }

        mounted(): void {
            this.chart = new BaseChart(this.$refs.chart as HTMLCanvasElement);

            this.findDataset().then(dataset => {
                this.chart.init(this.chartType, dataset);
                this.flash();
            });
        }

        findDataset(): Promise<IDataset> {
            if (_.isUndefined(this.connection)) {
                throw new Error('Not connected.');
            }

            this.flash({
                message: 'Loading...',
                type: FlashType.Loading,
            });

            if (this.selected.name === 'Memory') {
                return this.connection.repository.findMemoryDataset(this.chartType);
            }

            if (this.selected.name === 'Disk') {
                return this.connection.repository.findDiskDataset(
                    this.chartType,
                    this.selected.path as string,
                );
            }

            return this.connection.repository.findCpuDataset(this.chartType);
        }
    };
</script>
