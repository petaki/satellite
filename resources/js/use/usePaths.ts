export default () => {
    const probePath = () => '/';

    const cpuPath = (probe: string, type?: string) => {
        const params = new URLSearchParams({ probe });

        if (type) {
            params.set('type', type);
        }

        return `/cpu?${params.toString()}`;
    };

    const memoryPath = (probe: string, type?: string) => {
        const params = new URLSearchParams({ probe });

        if (type) {
            params.set('type', type);
        }

        return `/memory?${params.toString()}`;
    };

    const loadPath = (probe: string, type?: string) => {
        const params = new URLSearchParams({ probe });

        if (type) {
            params.set('type', type);
        }

        return `/load?${params.toString()}`;
    };

    const diskPath = (probe: string, path?: string, type?: string) => {
        const params = new URLSearchParams({ probe });

        if (path) {
            params.set('path', path);
        }

        if (type) {
            params.set('type', type);
        }

        return `/disk?${params.toString()}`;
    };

    const logPath = (probe: string, path?: string) => {
        const params = new URLSearchParams({ probe });

        if (path) {
            params.set('path', path);
        }

        return `/log?${params.toString()}`;
    };

    const probeDeletePath = (probe: string) => {
        const params = new URLSearchParams({ probe });

        return `/probe/delete?${params.toString()}`;
    };

    const probeDeleteAllPath = () => '/probe/delete-all';

    return {
        cpuPath,
        diskPath,
        loadPath,
        logPath,
        memoryPath,
        probePath,
        probeDeleteAllPath,
        probeDeletePath
    };
};
