<template>
    <button :class="modelValue ? 'btn-pill-green' : 'btn-pill-yellow'"
            @click="$emit('update:modelValue', !modelValue)">
        <span class="h-2 w-2 rounded-full"
              :class="modelValue
                  ? 'bg-green-500 dark:bg-green-400'
                  : 'bg-yellow-500 dark:bg-yellow-400'"></span>
        {{ modelValue ? 'Live' : 'Paused' }}
    </button>
</template>

<script setup lang="ts">
defineProps<{
    modelValue: boolean
}>();

defineEmits<{
    'update:modelValue': [value: boolean]
}>();
</script>
