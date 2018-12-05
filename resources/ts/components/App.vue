<template>
    <div class="flex">
        <sidebar></sidebar>
        <probe></probe>
    </div>
</template>

<script lang="ts">
    import { Component, Vue } from "vue-property-decorator";
    import Probe from "./Probe.vue";
    import Sidebar from "./Sidebar.vue";

    @Component({
        name: "App",
        
        components: {
            Probe,
            Sidebar,
        }
    })
    export default class App extends Vue {};
</script>
