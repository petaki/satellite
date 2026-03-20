<template>
    <!-- eslint-disable max-len -->
    <div class="flex flex-col items-center justify-center h-full text-slate-400 dark:text-slate-500">
        <!-- eslint-enable max-len -->
        <chart-bar-icon class="h-12 w-12 mb-3" />
        <span class="text-sm">No data available.</span>
    </div>
</template>

<script setup lang="ts">
import {
    ChartBarIcon
} from '@heroicons/vue/24/outline';
</script>
