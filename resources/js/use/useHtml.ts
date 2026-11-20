const entities: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;'
};

export default function escapeHtml(text: string): string {
    return text.replace(/[&<>"']/g, character => entities[character]);
}
