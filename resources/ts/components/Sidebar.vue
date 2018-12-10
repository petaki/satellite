<template>
    <div class="sidebar">
        <h1 v-if="isConnected" class="logo">
            <i class="fas fa-wifi"></i>
            <small>
                Probe
            </small>
        </h1>
        <h1 v-else class="logo">
            <i class="fas fa-rocket"></i>
            <small>
                Carrier
            </small>
        </h1>
        <ul v-if="isConnected" class="nav">
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('cpu')}"
                   @click.prevent="select({name: 'cpu'})">
                    <i class="fas fa-microchip"></i>
                    CPU
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('memory')}"
                   @click.prevent="select({name: 'memory'})">
                    <i class="fas fa-memory"></i>
                    Memory
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('drive-0')}"
                   @click.prevent="select({name: 'drive-0'})">
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
                   :class="{active: isSelected('drive-1')}"
                   @click.prevent="select({name: 'drive-1'})">
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
                   :class="{active: isSelected('new')}"
                   @click.prevent="select({name: 'new'})">
                    <i class="fas fa-plus-square"></i>
                    New
                </a>
            </li>
            <li v-for="(probe, index) in probes" :key="index">
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('probe', probe)}"
                   @click.prevent="select({name: 'probe', probe: probe})">
                    <i class="fas fa-wifi"></i>
                    {{ probe.name }}
                </a>
            </li>
        </ul>
    </div>
</template>

<script lang="ts">
    import Manager from '../connection';
    import { IProbe, ISelected } from '../store/types';
    import { Component, Vue } from 'vue-property-decorator';
    import { Mutation, State } from 'vuex-class';

    @Component({
        name: 'Sidebar',
    })
    export default class Sidebar extends Vue {
        @State selected!: ISelected;
        @State probes!: IProbe[];
        @Mutation select!: Function;

        get isConnected(): boolean {
            return Manager.isConnected;
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
