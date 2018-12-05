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
                   @click.prevent="select('cpu')">
                    <i class="fas fa-microchip"></i>
                    CPU
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('memory')}"
                   @click.prevent="select('memory')">
                    <i class="fas fa-memory"></i>
                    Memory
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('drive-0')}"
                   @click.prevent="select('drive-0')">
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
                   @click.prevent="select('drive-1')">
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
                   @click.prevent="select('new')">
                    <i class="fas fa-plus-square"></i>
                    New
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn"
                   :class="{active: isSelected('probe')}"
                   @click.prevent="select('probe')">
                    <i class="fas fa-wifi"></i>
                    Probe
                </a>
            </li>
        </ul>
    </div>
</template>

<script lang="ts">
    import { Component, Vue } from "vue-property-decorator";

    @Component({
        name: "Sidebar",
    })
    export default class Sidebar extends Vue {
        isConnected = false;
        selected = "new";

        isSelected(selected: string): boolean {
            return this.selected === selected;
        }

        select(selected: string): void {
            this.selected = selected;
        }
    };
</script>
