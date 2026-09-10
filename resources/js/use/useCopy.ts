import { ref } from 'vue';

export default (timeout = 2000) => {
    const copied = ref(false);
    let copiedTimeout: ReturnType<typeof setTimeout>;

    const copy = async (text: string) => {
        if (!navigator.clipboard) {
            return;
        }

        try {
            await navigator.clipboard.writeText(text);
        } catch {
            return;
        }

        copied.value = true;
        clearTimeout(copiedTimeout);
        copiedTimeout = setTimeout(() => {
            copied.value = false;
        }, timeout);
    };

    return {
        copied,
        copy
    };
};
