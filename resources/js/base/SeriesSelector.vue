<template>
    <toggle-button-group v-model="model"
                         :options="seriesButtons"
                         button-class="md:border-r-0 md:rounded-none hidden lg:block"
                         first-button-class="md:rounded-l-sm md:rounded-r-none" />
    <select v-model="model"
            class="form-select h-11 md:rounded-none">
        <option v-for="type in $page.props.seriesTypes"
                :key="type.value"
                :value="type.value">
            {{ type.name }}
        </option>
    </select>
</template>

<script setup lang="ts">
import {
    computed,
    watch
} from 'vue';

import { router, usePage } from '@inertiajs/vue3';
import type { SeriesType } from '../types';
import ToggleButtonGroup from './ToggleButtonGroup.vue';

const model = defineModel<string>();

const { href } = defineProps<{
    href:(isDefault: boolean, selectedType: string | undefined) => string
}>();

const page = usePage();

model.value = page.props.seriesType as string;

const seriesButtons = computed(() => (page.props.seriesTypes as SeriesType[])
    .filter(type => (page.props.seriesButtons as string[]).indexOf(type.value) !== -1)
    .map(type => ({
        label: type.name.split(' ').map(segment => (!Number.isNaN(Number(segment))
            ? segment
            : segment[0])).join('').toUpperCase(),
        value: type.value
    })));

const loadSeries = () => {
    router.visit(href(
        model.value === page.props.seriesType,
        model.value
    ));
};

watch(model, () => {
    loadSeries();
});
</script>
