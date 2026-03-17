<template>
    <inertia-head :title="title
        ? `${title} - ${$page.props.title}${$page.props.suffix}`
        : `${$page.props.title}${$page.props.suffix}`">
        <slot></slot>
    </inertia-head>
</template>

<script setup lang="ts">
defineProps({
    title: {
        type: String,
        default: ''
    }
});
</script>
