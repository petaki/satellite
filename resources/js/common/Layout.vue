<template>
    <section class="d-flex">
        <div v-show="isSidebarOpen"
             class="sidebar__backdrop"
             @click.prevent="isSidebarOpen = false"></div>
        <div class="sidebar" :class="{'sidebar--open': isSidebarOpen}">
            <div class="sidebar__header">
                <inertia-link class="sidebar__logo" href="/">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('broadcast')" />
                    </svg>
                    <span>
                        {{ $page.props.title }}
                    </span>
                </inertia-link>
            </div>
            <div class="sidebar__content"></div>
            <div class="sidebar__footer">
                <span class="m-auto">
                    &copy; {{ year }}
                    <span class="sidebar__highlight">
                        {{ $page.props.title }}
                    </span>
                </span>
            </div>
        </div>
        <div class="content">
            <header class="content__header">
                <a class="content__header--toggler"
                   href="#"
                   @click.prevent="isSidebarOpen = true">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('list')" />
                    </svg>
                </a>
                <flash-message />
            </header>
            <main class="content__main">
                <slot></slot>
            </main>
        </div>
    </section>
</template>

<script>
import { Inertia } from '@inertiajs/inertia';
import FlashMessage from './FlashMessage.vue';

export default {
    components: {
        FlashMessage
    },

    data() {
        return {
            isSidebarOpen: false,
            year: new Date().getFullYear()
        };
    },

    mounted() {
        this.$once(
            'hook:destroyed',
            Inertia.on('navigate', () => { this.isSidebarOpen = false; })
        );
    }
};
</script>
