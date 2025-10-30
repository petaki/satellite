<template>
    <!-- eslint-disable max-len -->
    <inertia-link class="flex items-center px-5 py-3.5 w-full hover:bg-black/20"
                  :class="{
                      'text-cyan-500 bg-black/20': isActive,
                      'text-slate-300 hover:text-cyan-500': !isHighlight,
                      'text-red-400 hover:text-red-500': isHighlight,
                  }"
                  :href="href"
                  :method="method">
        <slot></slot>
    </inertia-link>
    <!-- eslint-enable max-len -->
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
    isActive: {
        type: Boolean,
        default: false
    },

    isHighlight: {
        type: Boolean,
        default: false
    },

    href: {
        type: String,
        required: true
    },

    method: {
        type: String,
        default: 'get'
    }
});
</script>
