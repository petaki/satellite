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

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := app.seriesTypes()
	if !app.seriesTypeExists(seriesTypes, seriesType) {
		app.notFound(w)

		return
	}

	diskPaths, err := app.seriesRepository.FindDiskPaths()
	if err != nil {
		app.serverError(w, err)

		return
	}

	cpuSeries, err := app.seriesRepository.FindCpu(seriesType)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "cpu/Index", map[string]interface{}{
		"isCpuActive": true,
		"seriesType":  seriesType,
		"seriesTypes": seriesTypes,
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

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := app.seriesTypes()
	if !app.seriesTypeExists(seriesTypes, seriesType) {
		app.notFound(w)

		return
	}

	diskPaths, err := app.seriesRepository.FindDiskPaths()
	if err != nil {
		app.serverError(w, err)

		return
	}

	memorySeries, err := app.seriesRepository.FindMemory(seriesType)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "memory/Index", map[string]interface{}{
		"isMemoryActive": true,
		"seriesType":     seriesType,
		"seriesTypes":    seriesTypes,
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

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := app.seriesTypes()
	if !app.seriesTypeExists(seriesTypes, seriesType) {
		app.notFound(w)

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

	diskSeries, err := app.seriesRepository.FindDisk(seriesType, diskPath)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "disk/Index", map[string]interface{}{
		"seriesType":  seriesType,
		"seriesTypes": seriesTypes,
		"diskPath":    diskPath,
		"diskPaths":   diskPaths,
		"diskSeries":  diskSeries,
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

func (app *App) seriesTypeExists(seriesTypes []map[string]interface{}, seriesType models.SeriesType) bool {
	for _, current := range seriesTypes {
		if current["value"] == seriesType {
			return true
		}
	}

	return false
}

func (app *App) seriesTypes() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":  "Day",
			"value": models.Day,
		},
		{
			"name":  "Week",
			"value": models.Week,
		},
		{
			"name":  "Month",
			"value": models.Month,
		},
	}
}
