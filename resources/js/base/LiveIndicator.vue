<template>
    <!-- eslint-disable max-len -->
    <button class="flex items-center gap-1.5 px-3 py-1 text-xs rounded-xs border transition-colors"
            :class="modelValue
                ? 'border-green-500 text-green-600 hover:bg-green-50 dark:border-green-400 dark:text-green-400 dark:hover:bg-green-400/10'
                : 'border-yellow-500 text-yellow-600 hover:bg-yellow-50 dark:border-yellow-400 dark:text-yellow-400 dark:hover:bg-yellow-400/10'"
            @click="$emit('update:modelValue', !modelValue)">
        <span class="h-2 w-2 rounded-full"
              :class="modelValue ? 'bg-green-500 dark:bg-green-400' : 'bg-yellow-500 dark:bg-yellow-400'"></span>
        {{ modelValue ? 'Live' : 'Paused' }}
    </button>
    <!-- eslint-enable max-len -->
</template>

<script setup lang="ts">
defineProps<{
    modelValue: boolean
}>();

defineEmits<{
    'update:modelValue': [value: boolean]
}>();
</script>
