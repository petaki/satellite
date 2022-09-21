package web

import (
	"errors"
	"net/http"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) cpuIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)

		return
	}

	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := a.seriesTypes()
	if !a.seriesTypeExists(seriesTypes, seriesType) {
		a.notFound(w)

		return
	}

	diskPaths, err := a.seriesRepository.FindDiskPaths()
	if err != nil {
		a.serverError(w, err)

		return
	}

	cpuSeries, err := a.seriesRepository.FindCPU(seriesType)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var cpuAlarm float64 = 0

	alarm, err := a.alarmRepository.Find()
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		cpuAlarm = alarm.CPU
	}

	err = a.inertiaManager.Render(w, r, "cpu/Index", map[string]interface{}{
		"isCpuActive": true,
		"seriesType":  seriesType,
		"seriesTypes": seriesTypes,
		"diskPaths":   diskPaths,
		"cpuSeries":   cpuSeries,
		"cpuAlarm":    cpuAlarm,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) memoryIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := a.seriesTypes()
	if !a.seriesTypeExists(seriesTypes, seriesType) {
		a.notFound(w)

		return
	}

	diskPaths, err := a.seriesRepository.FindDiskPaths()
	if err != nil {
		a.serverError(w, err)

		return
	}

	memorySeries, err := a.seriesRepository.FindMemory(seriesType)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var memoryAlarm float64 = 0

	alarm, err := a.alarmRepository.Find()
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		memoryAlarm = alarm.Memory
	}

	err = a.inertiaManager.Render(w, r, "memory/Index", map[string]interface{}{
		"isMemoryActive": true,
		"seriesType":     seriesType,
		"seriesTypes":    seriesTypes,
		"diskPaths":      diskPaths,
		"memorySeries":   memorySeries,
		"memoryAlarm":    memoryAlarm,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) diskIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	seriesType := models.SeriesType(r.URL.Query().Get("type"))
	if seriesType == "" {
		seriesType = models.Day
	}

	seriesTypes := a.seriesTypes()
	if !a.seriesTypeExists(seriesTypes, seriesType) {
		a.notFound(w)

		return
	}

	diskPath := r.URL.Query().Get("path")
	if diskPath == "" {
		a.notFound(w)

		return
	}

	diskPaths, err := a.seriesRepository.FindDiskPaths()
	if err != nil {
		a.serverError(w, err)

		return
	}

	if !a.diskPathExists(diskPaths, diskPath) {
		a.notFound(w)

		return
	}

	diskSeries, err := a.seriesRepository.FindDisk(seriesType, diskPath)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var diskAlarm float64 = 0

	alarm, err := a.alarmRepository.Find()
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		diskAlarm = alarm.Disk
	}

	err = a.inertiaManager.Render(w, r, "disk/Index", map[string]interface{}{
		"seriesType":  seriesType,
		"seriesTypes": seriesTypes,
		"diskPath":    diskPath,
		"diskPaths":   diskPaths,
		"diskSeries":  diskSeries,
		"diskAlarm":   diskAlarm,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) diskPathExists(diskPaths []string, diskPath string) bool {
	for _, current := range diskPaths {
		if current == diskPath {
			return true
		}
	}

	return false
}

func (a *app) seriesTypeExists(seriesTypes []map[string]interface{}, seriesType models.SeriesType) bool {
	for _, current := range seriesTypes {
		if current["value"] == seriesType {
			return true
		}
	}

	return false
}

func (a *app) seriesTypes() []map[string]interface{} {
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
