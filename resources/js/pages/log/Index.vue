<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
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
            </card-title>
            <!-- eslint-disable vue/html-indent -->
            <pre v-if="logPath && logEntries && logEntries.length"
                 class="bg-gray-900 p-4 text-xs whitespace-pre-wrap text-slate-300"><template v-for="(entry, index) in logEntries" :key="entry.timestamp"><template v-if="index > 0">

</template><span class="font-bold text-cyan-400">// --- {{ formatTimestamp(entry.timestamp) }} ---</span>
{{ entry.content || 'No output.' }}</template></pre>
            <!-- eslint-enable vue/html-indent -->
            <pre v-else-if="logPath"
                 class="bg-gray-900 p-4 text-xs whitespace-pre-wrap text-slate-300"><code>No log entries found.</code></pre>
        </div>
        <!-- eslint-enable max-len -->
    </div>
</template>

<script setup lang="ts">
import {
    DocumentTextIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    onMounted,
    onUnmounted
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';
import useDate from '../../use/useDate';

interface LogEntry {
    timestamp: number
    content: string
}

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

const subtitle = ref('Log');
const reloadTimer = 60000;
let reloadInterval: ReturnType<typeof setInterval>;

const { timestamp: formatTimestamp } = useDate();

const links = ref([
    { name: subtitle }
]);

const onPathChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;

    router.visit(`/log?probe=${probe}&path=${target.value}`);
};

onMounted(() => {
    reloadInterval = setInterval(() => { router.reload(); }, reloadTimer);
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});
</script>
