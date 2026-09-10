import type {
    LogEntry, GroupedSnapshot, LineDiff, LogStatus
} from '../types';

export function computeDiff(prevLines: string[], currentLines: string[]): LineDiff {
    const prevSet = new Set(prevLines);
    const currentSet = new Set(currentLines);

    const added: string[] = [];
    const removed: string[] = [];
    let unchangedCount = 0;

    currentLines.forEach(line => {
        if (!prevSet.has(line)) {
            added.push(line);
        } else {
            unchangedCount += 1;
        }
    });

    prevLines.forEach(line => {
        if (!currentSet.has(line)) {
            removed.push(line);
        }
    });

    return { added, removed, unchangedCount };
}

const errorPattern = /\b(?:fatal|error(?:s|ed|ing)?|panic(?:s|ked|king)?|exceptions?|fail(?:s|ed|ing|ure|ures)?)\b/;
const warningPattern = /\b(?:warn(?:s|ed|ing|ings)?|deprecat(?:ed|ing|ion|ions))\b/;

export function detectStatus(
    lines: string[],
    diff: LineDiff | undefined,
    content: string,
    prevContent: string | null
): LogStatus {
    const joined = lines.join('\n').toLowerCase();

    if (errorPattern.test(joined)) {
        return 'error';
    }

    if (warningPattern.test(joined)) {
        return 'warning';
    }

    if (prevContent !== null && content === prevContent) {
        return 'unchanged';
    }

    if (diff && diff.added.length + diff.removed.length <= 3) {
        return 'minor';
    }

    return 'changed';
}

export function groupSnapshots(entries: LogEntry[]): GroupedSnapshot[] {
    if (!entries.length) {
        return [];
    }

    const snapshots = entries.map(entry => ({
        entry,
        content: (entry.content || '').trim(),
        lines: (entry.content || '').split('\n')
    }));

    const groups: GroupedSnapshot[] = [];
    let i = 0;

    while (i < snapshots.length) {
        const current = snapshots[i];
        let endIndex = i;

        while (
            endIndex + 1 < snapshots.length
            && snapshots[endIndex + 1].content === current.content
        ) {
            endIndex += 1;
        }

        const repeatCount = endIndex - i + 1;

        // Chronological predecessor is the next item in the array (entries are newest-first)
        const chronologicalPredecessor = endIndex + 1 < snapshots.length
            ? snapshots[endIndex + 1]
            : null;

        const diff = chronologicalPredecessor
            ? computeDiff(chronologicalPredecessor.lines, current.lines)
            : undefined;

        const prevContent = chronologicalPredecessor
            ? chronologicalPredecessor.content
            : null;

        const status = detectStatus(current.lines, diff, current.content, prevContent);

        const preview = current.lines.find(l => l.trim().length > 0) || '';

        groups.push({
            id: String(current.entry.timestamp),
            startMinute: snapshots[endIndex].entry.timestamp,
            endMinute: current.entry.timestamp,
            repeatCount,
            lineCount: current.lines.length,
            status,
            preview: preview.length > 80 ? `${preview.slice(0, 80)}...` : preview,
            lines: current.lines,
            diff
        });

        i = endIndex + 1;
    }

    return groups;
}

export function filterSnapshots(
    groups: GroupedSnapshot[],
    query: string
): GroupedSnapshot[] {
    let result = groups;

    if (query.trim()) {
        const lower = query.toLowerCase();

        result = result.filter(
            g => g.lines.some(line => line.toLowerCase().includes(lower))
        );
    }

    return result;
}

export function generateSummary(group: GroupedSnapshot): string {
    const parts: string[] = [];

    parts.push(`${group.lineCount} lines, status: ${group.status}`);

    if (group.diff) {
        parts.push(`Changes: +${group.diff.added.length} / -${group.diff.removed.length} lines`);
    }

    if (group.repeatCount > 1) {
        parts.push(`Repeated ${group.repeatCount} times`);
    }

    return parts.join('\n');
}
