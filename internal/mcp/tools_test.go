package mcp

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/petaki/satellite/internal/models"
)

type probeRepository struct {
	models.ProbeRepository
}

func (probeRepository) FindAll() ([]models.Probe, error) {
	return []models.Probe{"web-01"}, nil
}

type logRepository struct{}

func (logRepository) FindLogPaths(models.Probe) ([]string, error) {
	return []string{"/var/log/syslog"}, nil
}

func (logRepository) FindLog(models.Probe, string) ([]models.LogEntry, error) {
	entries := make([]models.LogEntry, 15)

	for i := range entries {
		entries[i] = models.LogEntry{Timestamp: int64(len(entries) - i), Content: strconv.Itoa(i)}
	}

	return entries, nil
}

func TestGetLogsLimitsEntries(t *testing.T) {
	h := &handler{probeRepository: probeRepository{}, logRepository: logRepository{}}

	tests := []struct {
		name    string
		limit   any
		want    int
		isError bool
	}{
		{name: "default", want: defaultLogLimit},
		{name: "explicit", limit: 3.0, want: 3},
		{name: "above available", limit: 100.0, want: 15},
		{name: "zero", limit: 0.0, isError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]any{"probe": "web-01"}
			if tt.limit != nil {
				args["limit"] = tt.limit
			}

			request := mcp.CallToolRequest{}
			request.Params.Arguments = args

			result, err := h.getLogs(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}

			if result.IsError != tt.isError {
				t.Fatalf("IsError = %v, want %v", result.IsError, tt.isError)
			}

			if tt.isError {
				return
			}

			var got struct {
				LogEntries []models.LogEntry `json:"logEntries"`
			}

			err = json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &got)
			if err != nil {
				t.Fatal(err)
			}

			if len(got.LogEntries) != tt.want {
				t.Errorf("len(logEntries) = %d, want %d", len(got.LogEntries), tt.want)
			}

			if got.LogEntries[0].Content != "0" {
				t.Errorf("logEntries[0] = %q, want newest entry", got.LogEntries[0].Content)
			}
		})
	}
}
