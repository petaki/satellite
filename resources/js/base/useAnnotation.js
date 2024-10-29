export default () => {
    const alarm = value => ({
        y: value,
        borderColor: '#ef4444',
        label: {
            borderColor: '#ef4444',
            style: {
                color: '#fff',
                fontSize: '0.75rem',
                fontWeight: 700,
                background: '#ef4444'
            },
            text: `Alarm: ${value}%`,
            offsetY: 5
        }
    });

    const max = value => ({
        y: value,
        borderColor: '#f6d757',
        label: {
            borderColor: '#f6d757',
            style: {
                color: '#1f2937',
                fontSize: '0.75rem',
                fontWeight: 700,
                background: '#f6d757'
            },
            text: `Max: ${value.toFixed(2)}%`,
            offsetY: 5
        }
    });

    return {
        alarm,
        max
    };
};
