package models

// LogEntry type.
type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

// LogRepository type.
type LogRepository interface {
	FindLogPaths(Probe) ([]string, error)
	FindLog(Probe, string) ([]LogEntry, error)
}
