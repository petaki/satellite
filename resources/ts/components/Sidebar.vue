<template>
    <div class="sidebar">
        <h1 v-if="hasConnection" class="logo">
            <i class="fas fa-wifi"></i>
            <small>
                {{ connection.probe.name }}
            </small>
        </h1>
        <h1 v-else class="logo">
            <i class="fas fa-rocket"></i>
            <small>
                Carrier
            </small>
        </h1>
        <ul v-if="hasConnection" class="nav">
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('CPU')}"
                   @click.prevent="select({name: 'CPU'})">
                    <i class="fas fa-microchip"></i>
                    CPU
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('Memory')}"
                   @click.prevent="select({name: 'Memory'})">
                    <i class="fas fa-memory"></i>
                    Memory
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('Drive-0')}"
                   @click.prevent="select({name: 'Drive-0'})">
                    <i class="fas fa-hdd"></i>
                    Drive
                    <small>
                        /mnt/hdd0/documents
                    </small>
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('Drive-1')}"
                   @click.prevent="select({name: 'Drive-1'})">
                    <i class="fas fa-hdd"></i>
                    Drive
                    <small>
                        /mnt/hdd1/pictures
                    </small>
                </a>
            </li>
        </ul>
        <ul v-else class="nav">
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('New')}"
                   @click.prevent="select({name: 'New'})">
                    <i class="fas fa-plus-square"></i>
                    New
                </a>
            </li>
            <li v-for="(probe, index) in probes" :key="index">
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('Edit', probe)}"
                   @click.prevent="select({name: 'Edit', probe: probe})"
                   @dblclick="connect(probe)">
                    <i class="fas fa-wifi"></i>
                    {{ probe.name }}
                </a>
            </li>
        </ul>
    </div>
</template>

<script lang="ts">
    import _ from 'lodash';
    import { IConnection, IProbe, ISelected } from '../store/types';
    import { Component, Vue } from 'vue-property-decorator';
    import { Action, Mutation, State } from 'vuex-class';

    @Component({
        name: 'Sidebar',
    })
    export default class Sidebar extends Vue {
        @State connection?: IConnection;
        @State selected!: ISelected;
        @State probes!: IProbe[];

        @Mutation select!: (selected: ISelected) => void;
        @Action connect!: (probe: IProbe) => void;

        get hasConnection(): boolean {
            return !_.isUndefined(this.connection);
        }

        isSelected(name: string, probe?: IProbe): boolean {
            let isSelected = this.selected.name === name;

            if (isSelected && probe) {
                isSelected = isSelected && this.selected.probe === probe; 
            }

            return isSelected;
        }
    };
</script>
