package mcp

import (
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/cli"
)

// NewHandler creates an MCP HTTP handler supporting both SSE and Streamable HTTP transports.
func NewHandler(
	cliApp *cli.App,
	appConfig *config.Config,
	probeRepository models.ProbeRepository,
	alarmRepository models.AlarmRepository,
	seriesRepository models.SeriesRepository,
	logRepository models.LogRepository,
) http.Handler {
	h := &handler{
		appConfig:        appConfig,
		probeRepository:  probeRepository,
		alarmRepository:  alarmRepository,
		seriesRepository: seriesRepository,
		logRepository:    logRepository,
	}

	mcpServer := server.NewMCPServer(
		cliApp.Name,
		cliApp.Version,
	)

	mcpServer.AddTool(
		mcp.NewTool("list_probes",
			mcp.WithDescription("List all monitored probes with their current status, latest CPU/memory/load values, and alarm thresholds"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
		),
		h.listProbes,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_cpu",
			mcp.WithDescription("Get CPU usage time series data for a probe, including min/max/avg and top 3 processes"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
			mcp.WithString("series", mcp.Description("Time range (e.g. last_5_minutes, last_1_hour, last_24_hours)")),
		),
		h.getCPU,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_memory",
			mcp.WithDescription("Get memory usage time series data for a probe, including min/max/avg and top 3 processes"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
			mcp.WithString("series", mcp.Description("Time range (e.g. last_5_minutes, last_1_hour, last_24_hours)")),
		),
		h.getMemory,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_load",
			mcp.WithDescription("Get load average time series data for a probe (1m, 5m, 15m averages)"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
			mcp.WithString("series", mcp.Description("Time range (e.g. last_5_minutes, last_1_hour, last_24_hours)")),
		),
		h.getLoad,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_disk",
			mcp.WithDescription("Get disk usage time series data for a probe, including min/max/avg and available paths"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
			mcp.WithString("series", mcp.Description("Time range (e.g. last_5_minutes, last_1_hour, last_24_hours)")),
			mcp.WithString("path", mcp.Description("Disk path (defaults to first available path)")),
		),
		h.getDisk,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_logs",
			mcp.WithDescription("Get log entries for a probe, with available log file paths"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
			mcp.WithString("path", mcp.Description("Log file path (defaults to first available path)")),
		),
		h.getLogs,
	)

	mcpServer.AddTool(
		mcp.NewTool("get_alerts",
			mcp.WithDescription("Get alarm thresholds and current metric values for a probe"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
		),
		h.getAlerts,
	)

	mcpServer.AddTool(
		mcp.NewTool("delete_probe",
			mcp.WithDescription("Delete a probe and all its data"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithString("probe", mcp.Required(), mcp.Description("Probe name")),
		),
		h.deleteProbe,
	)

	httpServer := server.NewStreamableHTTPServer(mcpServer)

	return httpServer
}
