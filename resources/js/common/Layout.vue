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
            <div class="sidebar__content">
                <h1 class="sidebar__content--title">
                    Probe
                </h1>
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isPerformanceActive}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('cpu')" />
                            </svg>
                            <span>
                                Performance
                            </span>
                        </inertia-link>
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isDisksActive}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('hdd')" />
                            </svg>
                            <span>
                                Disks
                            </span>
                        </inertia-link>
                    </li>
                </ul>
            </div>
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
            </header>
            <main class="content__main">
                <slot></slot>
            </main>
        </div>
    </section>
</template>

<script>
import { Inertia } from '@inertiajs/inertia';

export default {
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
