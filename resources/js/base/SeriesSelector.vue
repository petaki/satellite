<template>
    <!-- eslint-disable max-len -->
    <button v-for="(button, index) in createSeriesButtons()"
            :key="button.value"
            class="btn-white py-0 px-3 h-11 md:border-r-0 md:rounded-none text-sm font-semibold"
            :class="{'md:rounded-l-sm md:rounded-r-none': index === 0, 'bg-gray-100': model === button.value}"
            type="button"
            @click="model = button.value">
        {{ button.name }}
    </button>
    <!-- eslint-enable max-len -->
    <select v-model="model"
            class="form-select h-11 border-gray-300 md:rounded-none">
        <option v-for="type in $page.props.seriesTypes"
                :key="type.value"
                :value="type.value">
            {{ type.name }}
        </option>
    </select>
    <button class="btn-white py-0 px-3 h-11 md:border-l-0 md:rounded-l-none"
            @click="loadSeries()">
        <arrow-path-icon class="h-5 w-5" />
    </button>
</template>

<script setup>
import {
    defineModel,
    defineProps,
    watch
} from 'vue';

import { ArrowPathIcon } from '@heroicons/vue/20/solid';
import { router, usePage } from '@inertiajs/vue3';

const model = defineModel({
    type: String
});

const { href } = defineProps({
    href: {
        type: Function,
        required: true
    }
});

const page = usePage();

model.value = page.props.seriesType;

const createSeriesButtons = () => page.props.seriesTypes
    .filter(type => page.props.seriesButtons.indexOf(type.value) !== -1)
    .map(type => ({
        name: type.name.match(/\b\w/g).join('').toUpperCase(),
        value: type.value
    }));

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
