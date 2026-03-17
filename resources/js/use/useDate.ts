import { DateTime, Duration } from 'luxon';

export default () => {
    const date = value => {
        const d = DateTime.fromISO(value);

        if (!d.isValid) {
            return value;
        }

        return d.toLocaleString(DateTime.DATETIME_MED);
    };

    const duration = minutes => Duration.fromObject({ minutes }).toFormat('hh:mm:ss');

    return {
        date,
        duration
    };
};
