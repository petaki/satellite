<template>
    <section class="flex">
        <div v-show="isSidebarOpen"
             class="bg-black bg-opacity-50 fixed inset-0 z-10 md:hidden"
             @click.prevent="isSidebarOpen = false"></div>
        <!-- eslint-disable max-len -->
        <div class="fixed inset-y-0 left-0 bg-gray-800 w-60 z-20 transform transition-transform md:static md:translate-x-0"
             :class="{'translate-x-0': isSidebarOpen, '-translate-x-full': !isSidebarOpen}">
            <div class="flex h-20 bg-gray-900">
                <inertia-link class="flex items-center m-auto text-white text-xl"
                              href="/">
                    <paper-airplane-icon class="h-7 w-7 mr-2" />
                    <span>
                        {{ $page.props.title }}
                    </span>
                </inertia-link>
            </div>
            <div class="sidebar overflow-y-auto">
                <sidebar-title>
                    Performance
                </sidebar-title>
                <sidebar-link :is-active="!!$page.props.isCpuActive"
                              href="/">
                    <chip-icon class="h-5 w-5 mr-2" />
                    <span>
                        CPU
                    </span>
                </sidebar-link>
                <sidebar-link :is-active="!!$page.props.isMemoryActive"
                              href="/memory">
                    <duplicate-icon class="h-5 w-5 mr-2" />
                    <span>
                        Memory
                    </span>
                </sidebar-link>
                <sidebar-title>
                    Disks
                </sidebar-title>
                <div v-if="!$page.props.diskPaths" class="px-5 py-3.5 text-blueGray-500">
                    No data found.
                </div>
                <sidebar-link v-for="diskPath in $page.props.diskPaths"
                              :key="diskPath"
                              :is-active="$page.props.diskPath === diskPath"
                              :href="`/disk?path=${diskPath}`">
                    <database-icon class="h-5 w-5 mr-2" />
                    <div>
                        Disk
                        <div class="text-xs">
                            {{ diskPath }}
                        </div>
                    </div>
                </sidebar-link>
            </div>
            <div class="flex h-20 bg-blueGray-700 bg-opacity-40 text-sm text-blueGray-300">
                <span class="m-auto">
                    &copy; {{ year }}
                    <span class="text-cyan-500">
                        {{ $page.props.title }}
                    </span>
                </span>
            </div>
        </div>
        <div class="flex-1">
            <header class="flex items-center bg-white h-20 shadow-sm px-5">
                <a class="md:hidden"
                   href="#"
                   @click.prevent="isSidebarOpen = true">
                    <menu-icon class="h-6 w-6" />
                </a>
            </header>
            <main class="content overflow-y-auto">
                <slot></slot>
            </main>
        </div>
        <!-- eslint-enable max-len -->
    </section>
</template>

<script>
import {
    MenuIcon,
    ChipIcon,
    DatabaseIcon,
    DuplicateIcon,
    PaperAirplaneIcon
} from '@heroicons/vue/outline';

import { ref, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import SidebarTitle from './SidebarTitle.vue';
import SidebarLink from './SidebarLink.vue';

export default {
    components: {
        MenuIcon,
        ChipIcon,
        DatabaseIcon,
        DuplicateIcon,
        PaperAirplaneIcon,
        SidebarTitle,
        SidebarLink
    },

    setup() {
        const isSidebarOpen = ref(false);
        const year = ref(new Date().getFullYear());

        onUnmounted(
            Inertia.on('navigate', () => {
                isSidebarOpen.value = false;
            })
        );

        return {
            isSidebarOpen,
            year
        };
    }
};
</script>
