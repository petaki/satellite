<template>
    <div id="app" class="flex">
        <sidebar/>
    </div>
</template>

<script lang="ts">
    import { Component, Vue } from "vue-property-decorator";
    import Sidebar from "./Sidebar.vue";

    @Component({
        name: "App",
        
        components: {
            Sidebar,
        }
    })
    export default class App extends Vue {};
</script>
