<template>
    <inertia-head :title="title ? `${title} - ${$page.props.title}` : $page.props.title">
        <slot></slot>
    </inertia-head>
</template>

<script>
export default {
    props: {
        title: {
            type: String,
            default: ''
        }
    }
};
</script>
