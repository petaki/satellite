<template>
    <div class="flex mb-5">
        <inertia-link class="text-cyan-500 hover:text-cyan-600"
                      href="/">
            Home
        </inertia-link>
        <template v-for="(link, index) in links"
                  :key="link.name">
            <div v-if="index === 0"
                 class="px-2 text-slate-300">
                /
            </div>
            <inertia-link v-if="link.href"
                          class="link"
                          :href="link.href">
                {{ link.name }}
            </inertia-link>
            <div v-else class="dark:text-slate-300">
                {{ link.name }}
            </div>
            <div v-if="index < links.length - 1"
                 class="px-2 text-slate-300">
                /
            </div>
        </template>
    </div>
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
    links: {
        type: Array,
        required: true
    }
});
</script>
