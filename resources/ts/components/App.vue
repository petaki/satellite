<template>
    <div class="flex">
        <sidebar></sidebar>
        <probe v-if="isProbe"></probe>
        <chart v-else></chart>
        <flash></flash>
    </div>
</template>

<script lang="ts">
    import _ from 'lodash';
    import Chart from './Chart.vue';
    import Flash from './Flash.vue';
    import Probe from './Probe.vue';
    import Sidebar from './Sidebar.vue';
    import { Component, Vue } from 'vue-property-decorator';
    import { ISelected } from '../store/types';
    import { State } from 'vuex-class';

    @Component({
        name: 'App',
        
        components: {
            Chart,
            Flash,
            Probe,
            Sidebar,
        }
    })
    export default class App extends Vue {
        @State selected!: ISelected;

        get isProbe(): boolean {
            return _.includes(['New', 'Edit'], this.selected.name);
        }
    };
</script>
