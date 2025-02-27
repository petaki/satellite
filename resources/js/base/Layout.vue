<template>
    <div v-show="isSidebarOpen"
         class="bg-black bg-opacity-50 fixed inset-0 z-20 md:hidden"
         @click.prevent="isSidebarOpen = false"></div>
    <!-- eslint-disable max-len -->
    <div class="fixed inset-y-0 left-0 bg-gray-800 w-60 z-30 transform transition-transform md:translate-x-0"
         :class="{'translate-x-0': isSidebarOpen, '-translate-x-full': !isSidebarOpen}">
        <div class="flex h-20 bg-gray-900">
            <inertia-link class="flex items-center m-auto text-white text-xl"
                          href="/">
                <paper-airplane-icon class="h-7 w-7 mr-2" />
                <span class="break-all">
                    {{ $page.props.title }}
                </span>
            </inertia-link>
        </div>
        <div class="sidebar overflow-y-auto">
            <div class="pb-7">
                <sidebar-title>
                    Probe
                </sidebar-title>
                <sidebar-link :is-active="!!$page.props.isProbeActive"
                              href="/">
                    <cube-icon v-if="!!$page.props.probe" class="h-5 w-5 mr-2" />
                    <cube-transparent-icon v-else class="h-5 w-5 mr-2" />
                    <span v-if="!!$page.props.probe">
                        {{ $page.props.probe }}
                    </span>
                    <span v-else>
                        Not selected.
                    </span>
                    <chevron-right-icon class="h-5 w-5 ml-auto" />
                </sidebar-link>
            </div>
            <div v-if="!!$page.props.probe" class="bg-black/20 pb-7">
                <sidebar-title>
                    Performance
                </sidebar-title>
                <sidebar-link :is-active="!!$page.props.isCpuActive"
                              :href="`/cpu?probe=${$page.props.probe}`">
                    <cpu-chip-icon class="h-5 w-5 mr-2" />
                    <span>
                        CPU
                    </span>
                </sidebar-link>
                <sidebar-link :is-active="!!$page.props.isMemoryActive"
                              :href="`/memory?probe=${$page.props.probe}`">
                    <document-duplicate-icon class="h-5 w-5 mr-2" />
                    <span>
                        Memory
                    </span>
                </sidebar-link>
                <sidebar-link :is-active="!!$page.props.isLoadActive"
                              :href="`/load?probe=${$page.props.probe}`">
                    <chart-bar-icon class="h-5 w-5 mr-2" />
                    <span>
                        Load
                    </span>
                </sidebar-link>
                <sidebar-title>
                    Disks
                </sidebar-title>
                <div v-if="!$page.props.diskPaths" class="px-5 py-3.5 text-slate-500">
                    No data found.
                </div>
                <sidebar-link v-for="diskPath in $page.props.diskPaths"
                              :key="diskPath"
                              :is-active="$page.props.diskPath === diskPath"
                              :href="`/disk?probe=${$page.props.probe}&path=${diskPath}`">
                    <circle-stack-icon class="h-5 w-5 mr-2" />
                    <div class="flex-1">
                        Disk
                        <div class="text-xs break-all">
                            {{ diskPath }}
                        </div>
                    </div>
                </sidebar-link>
                <sidebar-title>
                    Delete
                </sidebar-title>
                <!-- eslint-disable vue/attribute-hyphenation -->
                <sidebar-link :is-highlight="true"
                              :href="`/probe/delete?probe=${$page.props.probe}`"
                              method="delete"
                              as="button"
                              :onBefore="confirmDelete">
                    <trash-icon class="h-5 w-5 mr-2" />
                    <span>
                        Delete
                    </span>
                </sidebar-link>
                <!-- eslint-enable vue/attribute-hyphenation -->
            </div>
        </div>
        <div class="flex h-20 bg-slate-700 bg-opacity-40 text-sm text-slate-300">
            <span class="m-auto">
                &copy; {{ year }}
                <span class="text-cyan-500">
                    {{ $page.props.title }}
                </span>
            </span>
        </div>
    </div>
    <div class="md:ml-60">
        <header class="flex items-center bg-white h-20 shadow-xs px-5">
            <a class="md:hidden"
               href="#"
               @click.prevent="isSidebarOpen = true">
                <bars3-icon class="h-6 w-6" />
            </a>
        </header>
        <main class="content overflow-y-auto" scroll-region>
            <slot></slot>
        </main>
    </div>
    <!-- eslint-enable max-len -->
</template>

<script setup>
import {
    Bars3Icon,
    CpuChipIcon,
    CircleStackIcon,
    ChartBarIcon,
    CubeIcon,
    CubeTransparentIcon,
    DocumentDuplicateIcon,
    PaperAirplaneIcon,
    TrashIcon
} from '@heroicons/vue/24/outline';

import {
    ChevronRightIcon
} from '@heroicons/vue/20/solid';

import { ref, onUnmounted } from 'vue';
import { router } from '@inertiajs/vue3';
import SidebarTitle from './SidebarTitle.vue';
import SidebarLink from './SidebarLink.vue';

const isSidebarOpen = ref(false);
const year = ref(new Date().getFullYear());

// eslint-disable-next-line no-alert
const confirmDelete = () => window.confirm('Are you sure you want to delete this probe?');

onUnmounted(
    router.on('navigate', () => {
        isSidebarOpen.value = false;
    })
);
</script>
