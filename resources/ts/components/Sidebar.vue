<template>
    <div class="w-64 bg-indigo-darker min-h-screen py-6">
        <h1 class="text-blue font-light text-center">
            <i class="fas fa-wifi"></i>
            <small class="block text-indigo text-sm mt-2">
                Probe
            </small>
        </h1>
        <ul class="list-reset font-light mt-6">
            <li>
                <a href="#"
                   class="btn-sidebar"
                   :class="{active: isSelected('cpu')}"
                   @click.prevent="select('cpu')">
                    <i class="fas fa-microchip"></i>
                    CPU
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn-sidebar"
                   :class="{active: isSelected('memory')}"
                   @click.prevent="select('memory')">
                    <i class="fas fa-memory"></i>
                    Memory
                </a>
            </li>
            <li>
                <a href="#"
                   class="btn-sidebar"
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
                   class="btn-sidebar"
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
    </div>
</template>

<script lang="ts">
    import { Component, Vue } from "vue-property-decorator";

    @Component({
        name: "Sidebar",
    })
    export default class Sidebar extends Vue {
        selected = "cpu";

        isSelected(selected: string): boolean {
            return this.selected === selected;
        }

        select(selected: string): void {
            this.selected = selected;
        }
    };
</script>
