import { DateTime, Duration } from 'luxon';

export default () => {
    const date = (value: string) => {
        const d = DateTime.fromISO(value);

        if (!d.isValid) {
            return value;
        }

        return d.toLocaleString(DateTime.DATETIME_MED);
    };

    const duration = (minutes: number) => Duration.fromObject({ minutes }).toFormat('hh:mm:ss');

    const timestamp = (value: number) => {
        const d = DateTime.fromSeconds(value);

        if (!d.isValid) {
            return String(value);
        }

        return d.toLocaleString(DateTime.DATETIME_MED_WITH_SECONDS);
    };

    return {
        date,
        duration,
        timestamp
    };
};
