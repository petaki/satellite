<template>
    <span><span v-for="(segment, i) in segments"
                :key="i"
                :style="segment.style">{{ segment.text }}</span></span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { parseAnsi } from '../use/useAnsi';

const { line } = defineProps<{
    line: string
}>();

const segments = computed(() => parseAnsi(line));
</script>
