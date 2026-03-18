package web

import (
	"net/http"
	"slices"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) logIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	logPaths, err := a.logRepository.FindLogPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	diskPaths, err := a.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	logPath := r.URL.Query().Get("path")

	if logPath == "" && len(logPaths) > 0 {
		logPath = logPaths[0]
	}

	var logEntries []models.LogEntry

	if logPath != "" && slices.Contains(logPaths, logPath) {
		logEntries, err = a.logRepository.FindLog(probe, logPath)
		if err != nil {
			a.serverError(w, err)

			return
		}
	} else {
		logPath = ""
	}

	err = a.inertiaManager.Render(w, r, "log/Index", map[string]any{
		"isLogActive": true,
		"logPath":     logPath,
		"logPaths":    logPaths,
		"logEntries":  logEntries,
		"diskPaths":   diskPaths,
	})
	if err != nil {
		a.serverError(w, err)
	}
}
