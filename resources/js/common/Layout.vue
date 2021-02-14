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
                    Performance
                </h1>
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isCpuActive}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('cpu')" />
                            </svg>
                            <span>
                                CPU
                            </span>
                        </inertia-link>
                    </li>
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isMemoryActive}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('grid-3x2')" />
                            </svg>
                            <span>
                                Memory
                            </span>
                        </inertia-link>
                    </li>
                </ul>
                <h1 class="sidebar__content--title">
                    Disks
                </h1>
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <inertia-link v-for="diskPath in $page.props.diskPaths"
                                      :key="diskPath"
                                      class="nav-link"
                                      :class="{active: $page.props.diskPath === diskPath}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('hdd')" />
                            </svg>
                            <span>
                                Disk
                                <span class="sidebar__muted">
                                    {{ diskPath }}
                                </span>
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