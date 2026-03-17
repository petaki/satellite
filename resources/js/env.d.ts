/// <reference types="vite/client" />

declare global {
    interface Window {
        isDark: () => boolean
        createApex: () => void
        Apex: any
    }
}

declare module '*.vue';

export {};
