<template>
    <!-- eslint-disable max-len -->
    <inertia-link class="flex items-center px-5 py-3.5 w-full hover:bg-black/20"
                  :class="{
                      'text-cyan-500 bg-black/20': isActive,
                      'text-slate-300 hover:text-cyan-500': !isHighlight,
                      'text-red-400 hover:text-red-500': isHighlight,
                  }"
                  :href="href"
                  :method="method"
                  :as="as"
                  :on-before="onBefore">
        <slot />
    </inertia-link>
    <!-- eslint-enable max-len -->
</template>

<script setup lang="ts">
import type { InertiaLinkProps } from '@inertiajs/vue3';

withDefaults(defineProps<{
    isActive?: boolean
    isHighlight?: boolean
    href: string
    method?: InertiaLinkProps['method']
    as?: InertiaLinkProps['as']
    onBefore?: InertiaLinkProps['onBefore']
}>(), {
    isActive: false,
    isHighlight: false,
    method: 'get',
    as: 'a',
    onBefore: () => {}
});
</script>
