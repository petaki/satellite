import { ref } from 'vue';

const execCopy = (text: string) => {
    const focused = document.activeElement;
    const textarea = document.createElement('textarea');

    textarea.value = text;
    textarea.setAttribute('readonly', '');
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();

    try {
        return document.execCommand('copy');
    } catch {
        return false;
    } finally {
        textarea.remove();

        if (focused instanceof HTMLElement) {
            focused.focus();
        }
    }
};

const writeText = async (text: string) => {
    if (!navigator.clipboard) {
        return execCopy(text);
    }

    try {
        await navigator.clipboard.writeText(text);
    } catch {
        return false;
    }

    return true;
};

export default (timeout = 2000) => {
    const copied = ref(false);
    let copiedTimeout: ReturnType<typeof setTimeout>;

    const copy = async (text: string) => {
        if (!await writeText(text)) {
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
