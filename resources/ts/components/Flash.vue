<template>
    <div v-if="hasFlash"
         class="absolute pin-t pin-l px-6 py-3 text-white"
         :class="{'bg-green-dark': isSuccess || isLoading, 'bg-red-light': isError}">
        <i v-if="isLoading" class="fas fa-circle-notch fa-spin mr-1"></i>
        {{ flash.message }}
    </div>
</template>
<script lang="ts">
    import _ from 'lodash';
    import { Component, Vue } from 'vue-property-decorator';
    import { IFlash, FlashType } from '../store/types';
    import { State } from 'vuex-class';

    @Component({
        name: 'Flash',
    })
    export default class Flash extends Vue {
        @State flash?: IFlash;

        get hasFlash(): boolean {
            return !_.isUndefined(this.flash);
        }

        get isSuccess(): boolean {
            return !_.isUndefined(this.flash) && this.flash.type === FlashType.Success;
        }

        get isLoading(): boolean {
            return !_.isUndefined(this.flash) && this.flash.type === FlashType.Loading;
        }

        get isError(): boolean {
            return !_.isUndefined(this.flash) && this.flash.type === FlashType.Error;
        }
    };
</script>
