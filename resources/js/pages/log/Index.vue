<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links">
            <live-indicator v-model="isLive" />
        </breadcrumb>
        <!-- eslint-disable max-len -->
        <div class="bg-white p-8 dark:bg-slate-700">
            <card-title>
                <document-text-icon class="h-6 w-6 shrink-0" />
                <select class="form-select flex-1 sm:ml-2"
                        :value="logPath || ''"
                        @change="onPathChange">
                    <option v-for="path in logPaths"
                            :key="path"
                            :value="path">
                        {{ path }}
                    </option>
                </select>
                <input v-if="logPath && logEntries && logEntries.length"
                       v-model="searchQuery"
                       class="form-input flex-1 sm:ml-2"
                       placeholder="Search">
                <template v-if="logPath && logEntries && logEntries.length">
                    <button class="btn-white py-0 px-3 h-11 sm:ml-2 md:rounded-r-none md:border-r-0 text-sm font-semibold"
                            type="button"
                            @click="cycleViewMode">
                        {{ viewModeLabel }}
                    </button>
                    <button class="btn-white py-0 px-3 h-11 md:rounded-l-none text-sm font-semibold"
                            type="button"
                            @click="jumpToLatest">
                        Latest
                    </button>
                </template>
            </card-title>
            <template v-if="logPath && logEntries && logEntries.length">
                <div class="flex gap-4 mt-4 chart">
                    <!-- Snapshot list (sidebar) -->
                    <div class="w-80 shrink-0 overflow-y-auto bg-gray-900">
                        <button v-for="group in filteredGroups"
                                :key="group.id"
                                class="w-full text-left px-3 py-2"
                                :class="group.id === selectedId
                                    ? 'bg-slate-600/30'
                                    : 'hover:bg-slate-700/30'"
                                @click="onSelect(group.id)">
                            <div class="flex items-center gap-2 text-xs">
                                <span class="text-slate-400 shrink-0">
                                    {{ formatRowTime(group) }}
                                </span>
                                <span class="text-slate-500 shrink-0">
                                    {{ group.lineCount }}L
                                    <template v-if="group.repeatCount > 1">
                                        x{{ group.repeatCount }}
                                    </template>
                                </span>
                                <log-status-badge class="ml-auto"
                                                  :status="group.status" />
                            </div>
                            <div class="text-xs text-slate-400 mt-1 truncate">
                                {{ group.preview || 'No output.' }}
                            </div>
                        </button>
                    </div>
                    <!-- Detail pane -->
                    <div class="flex-1 bg-gray-900 flex flex-col overflow-hidden">
                        <template v-if="selectedGroup">
                            <div class="flex items-center justify-between px-4 py-2
                                        border-b border-slate-700">
                                <div class="flex items-center gap-2 text-xs text-slate-400">
                                    <log-status-badge :status="selectedGroup.status" />
                                    <span>{{ detailTimeLabel }}</span>
                                    <span v-if="selectedGroup.repeatCount > 1">
                                        x{{ selectedGroup.repeatCount }}
                                    </span>
                                </div>
                                <button class="text-slate-400 hover:text-slate-200
                                               transition-colors p-1"
                                        :title="copied ? 'Copied!' : 'Copy to clipboard'"
                                        @click="copyContent">
                                    <clipboard-document-check-icon v-if="copied"
                                                                   class="h-4 w-4" />
                                    <clipboard-document-icon v-else
                                                             class="h-4 w-4" />
                                </button>
                            </div>
                            <div ref="contentEl"
                                 class="flex-1 overflow-y-auto p-4">
                                <pre v-if="viewMode === 'raw'"
                                     class="text-xs whitespace-pre-wrap text-slate-300">{{ selectedGroup.lines.join('\n') }}</pre>
                                <div v-else-if="viewMode === 'changes'"
                                     class="text-xs font-mono">
                                    <template v-if="selectedGroup.diff">
                                        <div v-for="(line, i) in selectedGroup.diff.added"
                                             :key="'a' + i"
                                             class="text-green-400">
                                            + {{ line }}
                                        </div>
                                        <div v-for="(line, i) in selectedGroup.diff.removed"
                                             :key="'r' + i"
                                             class="text-red-400">
                                            - {{ line }}
                                        </div>
                                        <div v-if="selectedGroup.diff.unchangedCount > 0"
                                             class="text-slate-500 mt-2 py-1 px-2
                                                    bg-slate-800 rounded text-center">
                                            {{ selectedGroup.diff.unchangedCount }} unchanged
                                            line{{ selectedGroup.diff.unchangedCount === 1
                                                ? '' : 's' }} hidden
                                        </div>
                                    </template>
                                    <div v-else class="text-slate-500 italic">
                                        First snapshot — no previous data to compare.
                                    </div>
                                </div>
                                <pre v-else
                                     class="text-xs whitespace-pre-wrap text-slate-300">{{ summary }}</pre>
                            </div>
                        </template>
                        <div v-else
                             class="flex-1 flex items-center justify-center
                                    text-slate-500 text-sm">
                            Select a snapshot to view details.
                        </div>
                    </div>
                </div>
            </template>
            <pre v-else-if="logPath"
                 class="bg-gray-900 p-4 text-xs whitespace-pre-wrap text-slate-300"><code>No log entries found.</code></pre>
        </div>
        <!-- eslint-enable max-len -->
    </div>
</template>

<script setup lang="ts">
import {
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
    DocumentTextIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    computed,
    watch,
    watchEffect,
    nextTick
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import LiveIndicator from '../../base/LiveIndicator.vue';
import LogStatusBadge from '../../base/LogStatusBadge.vue';
import type { GroupedSnapshot, LogEntry, LogViewMode } from '../../types';
import useCopy from '../../use/useCopy';
import useDate from '../../use/useDate';
import useLiveReload from '../../use/useLiveReload';
import {
    groupSnapshots, filterSnapshots, generateSummary
} from '../../use/useLogGrouping';

const {
    probe,
    logPath = '',
    logPaths = [],
    logEntries = []
} = defineProps<{
    probe: string
    logPath?: string
    logPaths?: string[]
    logEntries?: LogEntry[]
}>();

defineOptions({
    layout: Layout
});

const { copied, copy } = useCopy();
const { timestamp, shortTime } = useDate();
const { isLive } = useLiveReload();
const subtitle = ref('Log');

const links = ref([
    { name: subtitle }
]);

// UI state
const searchQuery = ref('');
const viewMode = ref<LogViewMode>('raw');
const selectedId = ref<string | null>(null);

// Refs
const contentEl = ref<HTMLElement | null>(null);

// Computed
const groups = ref<GroupedSnapshot[]>([]);

watchEffect(async () => {
    groups.value = await groupSnapshots(logEntries);
});

const filteredGroups = computed(() => filterSnapshots(
    groups.value,
    searchQuery.value
));

const selectedGroup = computed(() => groups.value.find(g => g.id === selectedId.value) || null);

const viewModeLabels: Record<LogViewMode, string> = {
    raw: 'Raw',
    changes: 'Changes',
    summary: 'Summary'
};

const viewModeLabel = computed(() => viewModeLabels[viewMode.value]);

const formatRowTime = (group: { repeatCount: number, startMinute: number, endMinute: number }) => {
    if (group.repeatCount > 1) {
        return `${shortTime(group.startMinute)}-${shortTime(group.endMinute)}`;
    }

    return shortTime(group.startMinute);
};

const detailTimeLabel = computed(() => {
    if (!selectedGroup.value) {
        return '';
    }

    if (selectedGroup.value.repeatCount > 1) {
        return `${timestamp(selectedGroup.value.startMinute)} - ${shortTime(selectedGroup.value.endMinute)}`;
    }

    return timestamp(selectedGroup.value.startMinute);
});

const summary = computed(() => {
    if (!selectedGroup.value) {
        return '';
    }

    return generateSummary(selectedGroup.value);
});

const cycleViewMode = () => {
    const modes: LogViewMode[] = ['raw', 'changes', 'summary'];
    const currentIndex = modes.indexOf(viewMode.value);

    viewMode.value = modes[(currentIndex + 1) % modes.length];
};

// Auto-select newest when live
watch([groups, isLive], () => {
    if (isLive.value && groups.value.length) {
        selectedId.value = groups.value[0].id;
    }
}, { immediate: true });

const onSelect = (id: string) => {
    selectedId.value = id;
};

const jumpToLatest = () => {
    if (groups.value.length) {
        selectedId.value = groups.value[0].id;
    }
};

const onPathChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;

    router.visit(`/log?probe=${probe}&path=${target.value}`);
};

const copyContent = () => {
    if (selectedGroup.value) {
        copy(selectedGroup.value.lines.join('\n'));
    }
};

// Scroll preservation for detail pane
let savedScrollTop = 0;
let userHasScrolled = false;

watch(() => selectedGroup.value?.id, async () => {
    userHasScrolled = false;

    await nextTick();

    if (contentEl.value) {
        contentEl.value.scrollTop = 0;
    }
});

watch(() => selectedGroup.value?.lines, async () => {
    if (!contentEl.value || !userHasScrolled) {
        return;
    }

    savedScrollTop = contentEl.value.scrollTop;

    await nextTick();

    contentEl.value.scrollTop = savedScrollTop;
});

watch(contentEl, el => {
    if (el) {
        el.addEventListener('scroll', () => {
            userHasScrolled = true;
        }, { passive: true });
    }
});
</script>
