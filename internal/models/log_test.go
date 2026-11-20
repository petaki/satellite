package models

import "testing"

func TestLogEntryPlain(t *testing.T) {
	entry := LogEntry{Timestamp: 1, Content: "#033[32mok#033[0m"}

	if got := entry.Plain(); got.Content != "ok" || got.Timestamp != 1 {
		t.Errorf("Plain() = %+v, want {Timestamp:1 Content:ok}", got)
	}

	if entry.Content != "#033[32mok#033[0m" {
		t.Errorf("Plain() changed the receiver: %q", entry.Content)
	}
}
