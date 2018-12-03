<template>
    <div class="w-64 bg-indigo-darker min-h-screen">
        <ul class="list-reset font-sans font-light mt-6">
            <li>
                <a class="btn-sidebar active" href="#">
                    <i class="fas fa-microchip"></i>
                    CPU
                </a>
            </li>
            <li>
                <a class="btn-sidebar" href="#">
                    <i class="fas fa-memory"></i>
                    Memory
                </a>
            </li>
            <li>
                <a class="btn-sidebar" href="#">
                    <i class="fas fa-hdd"></i>
                    Drive
                    <small>
                        /mnt/hdd0/documents
                    </small>
                </a>
            </li>
            <li>
                <a class="btn-sidebar" href="#">
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
    export default class Sidebar extends Vue {};
</script>
