<template>
    <button class="btn-pill-cyan"
            @click="router.reload()">
        <arrow-path-icon class="h-3 w-3" />
        Refresh
    </button>
</template>

<script setup lang="ts">
import { ArrowPathIcon } from '@heroicons/vue/20/solid';
import { router } from '@inertiajs/vue3';
</script>
