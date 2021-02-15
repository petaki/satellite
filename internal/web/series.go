package web

import (
	"net/http"

	"github.com/petaki/satellite/internal/models"
)

func (app *App) cpuIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)

		return
	}

	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	diskPaths, err := app.seriesRepository.FindDiskPaths()
	if err != nil {
		app.serverError(w, err)

		return
	}

	cpuSeries, err := app.seriesRepository.FindCpu(models.Day)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "cpu/Index", map[string]interface{}{
		"isCpuActive": true,
		"diskPaths":   diskPaths,
		"cpuSeries":   cpuSeries,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) memoryIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	diskPaths, err := app.seriesRepository.FindDiskPaths()
	if err != nil {
		app.serverError(w, err)

		return
	}

	memorySeries, err := app.seriesRepository.FindMemory(models.Day)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "memory/Index", map[string]interface{}{
		"isMemoryActive": true,
		"diskPaths":      diskPaths,
		"memorySeries":   memorySeries,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) diskIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	diskPath := r.URL.Query().Get("path")
	if diskPath == "" {
		app.notFound(w)

		return
	}

	diskPaths, err := app.seriesRepository.FindDiskPaths()
	if err != nil {
		app.serverError(w, err)

		return
	}

	if !app.diskPathExists(diskPaths, diskPath) {
		app.notFound(w)

		return
	}

	diskSeries, err := app.seriesRepository.FindDisk(models.Day, diskPath)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "disk/Index", map[string]interface{}{
		"diskPath":   diskPath,
		"diskPaths":  diskPaths,
		"diskSeries": diskSeries,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) diskPathExists(diskPaths []string, diskPath string) bool {
	for _, current := range diskPaths {
		if current == diskPath {
			return true
		}
	}

	return false
}
