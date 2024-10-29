<template>
    <inertia-head :title="title ? `${title} - ${$page.props.title}` : $page.props.title">
        <slot></slot>
    </inertia-head>
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
    title: {
        type: String,
        default: ''
    }
});
</script>
