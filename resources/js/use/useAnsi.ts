import type { AnsiSegment } from '../types';

const escape = String.fromCharCode(27);
const sequencePattern = new RegExp(`(?:${escape}|#033)\\[([0-9;?]*)([A-Za-z])`, 'g');

const tabPattern = /#011/g;

const palette = [
    '#64748b', '#f87171', '#4ade80', '#facc15', '#60a5fa', '#e879f9', '#22d3ee', '#e2e8f0',
    '#94a3b8', '#fca5a5', '#86efac', '#fde047', '#93c5fd', '#f0abfc', '#67e8f9', '#f8fafc'
];

interface AnsiState {
    color?: string
    background?: string
    bold: boolean
    dim: boolean
    italic: boolean
    underline: boolean
}

const emptyState = (): AnsiState => ({
    bold: false,
    dim: false,
    italic: false,
    underline: false
});

const cubeLevel = (value: number) => (value === 0 ? 0 : 55 + value * 40);

const color256 = (index: number): string | undefined => {
    if (!Number.isInteger(index) || index < 0 || index > 255) {
        return undefined;
    }

    if (index < 16) {
        return palette[index];
    }

    if (index < 232) {
        const value = index - 16;

        return `rgb(${cubeLevel(Math.floor(value / 36))}, ${cubeLevel(Math.floor(value / 6) % 6)}, ${cubeLevel(value % 6)})`;
    }

    const gray = 8 + (index - 232) * 10;

    return `rgb(${gray}, ${gray}, ${gray})`;
};

const extendedColor = (params: number[], at: number): [string | undefined, number] => {
    if (params[at] === 5) {
        return [color256(params[at + 1]), 2];
    }

    if (params[at] === 2) {
        const [r, g, b] = params.slice(at + 1, at + 4);

        if ([r, g, b].every(v => Number.isInteger(v) && v >= 0 && v <= 255)) {
            return [`rgb(${r}, ${g}, ${b})`, 4];
        }

        return [undefined, 4];
    }

    return [undefined, 0];
};

const applySgr = (state: AnsiState, raw: string): AnsiState => {
    const params = raw === '' ? [0] : raw.split(';').map(p => (p === '' ? 0 : Number(p)));
    let next = { ...state };

    for (let i = 0; i < params.length; i += 1) {
        const code = params[i];

        if (code === 0) {
            next = emptyState();
        } else if (code === 1) {
            next.bold = true;
        } else if (code === 2) {
            next.dim = true;
        } else if (code === 3) {
            next.italic = true;
        } else if (code === 4) {
            next.underline = true;
        } else if (code === 22) {
            next.bold = false;
            next.dim = false;
        } else if (code === 23) {
            next.italic = false;
        } else if (code === 24) {
            next.underline = false;
        } else if (code >= 30 && code <= 37) {
            next.color = palette[code - 30];
        } else if (code >= 90 && code <= 97) {
            next.color = palette[code - 90 + 8];
        } else if (code === 39) {
            next.color = undefined;
        } else if (code >= 40 && code <= 47) {
            next.background = palette[code - 40];
        } else if (code >= 100 && code <= 107) {
            next.background = palette[code - 100 + 8];
        } else if (code === 49) {
            next.background = undefined;
        } else if (code === 38 || code === 48) {
            const [color, used] = extendedColor(params, i + 1);

            if (code === 38) {
                next.color = color;
            } else {
                next.background = color;
            }

            i += used;
        }
    }

    return next;
};

const styleOf = (state: AnsiState): Record<string, string> => {
    const style: Record<string, string> = {};

    if (state.color) {
        style.color = state.color;
    }

    if (state.background) {
        style.backgroundColor = state.background;
    }

    if (state.bold) {
        style.fontWeight = '600';
    }

    if (state.dim) {
        style.opacity = '0.7';
    }

    if (state.italic) {
        style.fontStyle = 'italic';
    }

    if (state.underline) {
        style.textDecoration = 'underline';
    }

    return style;
};

export function stripAnsi(text: string): string {
    return text.replace(tabPattern, '\t').replace(sequencePattern, '');
}

export function parseAnsi(line: string): AnsiSegment[] {
    const text = line.replace(tabPattern, '\t');
    const segments: AnsiSegment[] = [];
    let state = emptyState();
    let last = 0;

    const push = (part: string) => {
        if (part !== '') {
            segments.push({ text: part, style: styleOf(state) });
        }
    };

    Array.from(text.matchAll(sequencePattern)).forEach(match => {
        push(text.slice(last, match.index));

        if (match[2] === 'm') {
            state = applySgr(state, match[1]);
        }

        last = match.index + match[0].length;
    });

    push(text.slice(last));

    return segments;
}
