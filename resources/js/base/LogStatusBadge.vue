<template>
    <span class="inline-flex items-center rounded-xs border px-2 py-0.5 text-xs font-medium"
          :class="badgeClass">
        {{ label }}
    </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import type { LogStatus } from '../types';

const { status } = defineProps<{
    status: LogStatus
}>();

const classes: Record<LogStatus, string> = {
    unchanged: 'border-gray-500 text-gray-400',
    changed: 'border-blue-500 text-blue-400',
    minor: 'border-sky-500 text-sky-400',
    warning: 'border-yellow-500 text-yellow-400',
    error: 'border-red-500 text-red-400'
};

const badgeClass = computed(() => classes[status]);

const label = computed(() => status.charAt(0).toUpperCase() + status.slice(1));
</script>
