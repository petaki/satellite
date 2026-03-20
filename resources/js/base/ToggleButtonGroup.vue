<template>
    <button v-for="(option, index) in options"
            :key="option.value"
            class="btn-white py-0 px-3 h-11 text-sm font-semibold"
            :class="[
                buttonClass,
                index === 0 ? firstButtonClass : '',
                index === options.length - 1 ? lastButtonClass : '',
                modelValue === option.value ? 'bg-gray-100 dark:bg-slate-800' : ''
            ]"
            type="button"
            @click="$emit('update:modelValue', option.value)">
        {{ option.label }}
    </button>
</template>

<script setup lang="ts">
defineProps<{
    options: { value: string, label: string }[]
    modelValue: string
    buttonClass?: string
    firstButtonClass?: string
    lastButtonClass?: string
}>();

defineEmits<{
    'update:modelValue': [value: string]
}>();
</script>
