import type {
    LogEntry, GroupedSnapshot, LineDiff, LogStatus
} from '../types';

export async function hashContent(content: string): Promise<string> {
    const data = new TextEncoder().encode(content.trim());
    const buffer = await crypto.subtle.digest('SHA-1', data);

    return Array.from(new Uint8Array(buffer))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('')
        .slice(0, 12);
}

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

export function detectStatus(
    lines: string[],
    diff: LineDiff | undefined,
    fingerprint: string,
    prevFingerprint: string | null
): LogStatus {
    const joined = lines.join('\n').toLowerCase();

    if (/\b(error|fatal|panic|exception|fail)\b/.test(joined)) {
        return 'error';
    }

    if (/\b(warn|warning|deprecated)\b/.test(joined)) {
        return 'warning';
    }

    if (prevFingerprint !== null && fingerprint === prevFingerprint) {
        return 'unchanged';
    }

    if (diff && diff.added.length + diff.removed.length <= 3) {
        return 'minor';
    }

    return 'changed';
}

export async function groupSnapshots(entries: LogEntry[]): Promise<GroupedSnapshot[]> {
    if (!entries.length) {
        return [];
    }

    const fingerprints = await Promise.all(
        entries.map(entry => hashContent(entry.content))
    );

    const withFingerprints = entries.map((entry, index) => ({
        entry,
        fingerprint: fingerprints[index],
        lines: (entry.content || '').split('\n')
    }));

    const groups: GroupedSnapshot[] = [];
    let i = 0;

    while (i < withFingerprints.length) {
        const current = withFingerprints[i];
        let endIndex = i;

        while (
            endIndex + 1 < withFingerprints.length
            && withFingerprints[endIndex + 1].fingerprint === current.fingerprint
        ) {
            endIndex += 1;
        }

        const repeatCount = endIndex - i + 1;

        // Chronological predecessor is the next item in the array (entries are newest-first)
        const chronologicalPredecessor = endIndex + 1 < withFingerprints.length
            ? withFingerprints[endIndex + 1]
            : null;

        const diff = chronologicalPredecessor
            ? computeDiff(chronologicalPredecessor.lines, current.lines)
            : undefined;

        const prevFingerprint = chronologicalPredecessor
            ? chronologicalPredecessor.fingerprint
            : null;

        const status = detectStatus(current.lines, diff, current.fingerprint, prevFingerprint);

        const preview = current.lines.find(l => l.trim().length > 0) || '';

        groups.push({
            id: `${current.entry.timestamp}-${current.fingerprint}`,
            startMinute: withFingerprints[endIndex].entry.timestamp,
            endMinute: current.entry.timestamp,
            repeatCount,
            lineCount: current.lines.length,
            status,
            preview: preview.length > 80 ? `${preview.slice(0, 80)}...` : preview,
            lines: current.lines,
            fingerprint: current.fingerprint,
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
