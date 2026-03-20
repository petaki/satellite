<template>
    <!-- eslint-disable max-len -->
    <span class="lg:ml-auto mr-2 flex items-center h-11 px-3 border border-gray-300 rounded-sm text-sm font-semibold text-cyan-500 dark:border-cyan-500">
        <clock-icon class="h-4 w-4 mr-1.5" />
        {{ duration(chunkSize) }}
    </span>
    <!-- eslint-enable max-len -->
</template>

<script setup lang="ts">
import { ClockIcon } from '@heroicons/vue/20/solid';
import useDate from '../use/useDate';

defineProps<{
    chunkSize: number
}>();

const { duration } = useDate();
</script>
