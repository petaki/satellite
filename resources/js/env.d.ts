/// <reference types="vite/client" />

import type { ApexConfig } from './types';

declare global {
    interface Window {
        isDark: () => boolean
        createApex: () => void
        Apex: ApexConfig
    }
}

declare module '*.vue';

export {};
