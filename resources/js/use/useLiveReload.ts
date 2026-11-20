import { ref, onMounted, onUnmounted } from 'vue';

import { router } from '@inertiajs/vue3';

export default (interval = 60000) => {
    const isLive = ref(true);
    let reloadInterval: ReturnType<typeof setInterval> | undefined;

    const reload = () => {
        if (isLive.value) {
            router.reload();
        }
    };

    const start = () => {
        reloadInterval = setInterval(reload, interval);
    };

    const stop = () => {
        clearInterval(reloadInterval);
    };

    const onVisibilityChange = () => {
        stop();

        if (document.hidden) {
            return;
        }

        reload();
        start();
    };

    onMounted(() => {
        if (!document.hidden) {
            start();
        }

        document.addEventListener('visibilitychange', onVisibilityChange);
    });

    onUnmounted(() => {
        stop();
        document.removeEventListener('visibilitychange', onVisibilityChange);
    });

    return {
        isLive
    };
};
