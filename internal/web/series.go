package web

import (
	"errors"
	"net/http"
	"slices"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) cpuIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	seriesType := a.seriesTypeFromRequest(r)

	if !a.seriesTypeExists(seriesType) {
		a.notFound(w)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeySeriesType, string(seriesType))

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	diskPaths, err := a.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	cpuMinSeries, cpuMaxSeries, cpuAvgSeries, process1Series, process2Series, process3Series, err := a.seriesRepository.FindCPU(probe, seriesType)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var cpuAlarm float64 = 0

	alarm, err := a.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		cpuAlarm = alarm.CPU
	}

	err = a.inertiaManager.Render(w, r, "cpu/Index", map[string]any{
		"isCpuActive":    true,
		"seriesType":     seriesType,
		"chunkSize":      a.seriesRepository.ChunkSize(seriesType),
		"diskPaths":      diskPaths,
		"cpuMinSeries":   cpuMinSeries,
		"cpuMaxSeries":   cpuMaxSeries,
		"cpuAvgSeries":   cpuAvgSeries,
		"process1Series": process1Series,
		"process2Series": process2Series,
		"process3Series": process3Series,
		"cpuAlarm":       cpuAlarm,
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

	seriesType := a.seriesTypeFromRequest(r)

	if !a.seriesTypeExists(seriesType) {
		a.notFound(w)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeySeriesType, string(seriesType))

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	diskPaths, err := a.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	memoryMinSeries, memoryMaxSeries, memoryAvgSeries, process1Series, process2Series, process3Series, err := a.seriesRepository.FindMemory(probe, seriesType)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var memoryAlarm float64 = 0

	alarm, err := a.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		memoryAlarm = alarm.Memory
	}

	err = a.inertiaManager.Render(w, r, "memory/Index", map[string]any{
		"isMemoryActive":  true,
		"seriesType":      seriesType,
		"chunkSize":       a.seriesRepository.ChunkSize(seriesType),
		"diskPaths":       diskPaths,
		"memoryMinSeries": memoryMinSeries,
		"memoryMaxSeries": memoryMaxSeries,
		"memoryAvgSeries": memoryAvgSeries,
		"process1Series":  process1Series,
		"process2Series":  process2Series,
		"process3Series":  process3Series,
		"memoryAlarm":     memoryAlarm,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) loadIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	seriesType := a.seriesTypeFromRequest(r)

	if !a.seriesTypeExists(seriesType) {
		a.notFound(w)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeySeriesType, string(seriesType))

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	diskPaths, err := a.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	load1Series, load5Series, load15Series, err := a.seriesRepository.FindLoad(probe, seriesType)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var loadAlarm float64 = 0

	alarm, err := a.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		loadAlarm = alarm.Load
	}

	err = a.inertiaManager.Render(w, r, "load/Index", map[string]any{
		"isLoadActive": true,
		"seriesType":   seriesType,
		"chunkSize":    a.seriesRepository.ChunkSize(seriesType),
		"diskPaths":    diskPaths,
		"load1Series":  load1Series,
		"load5Series":  load5Series,
		"load15Series": load15Series,
		"loadAlarm":    loadAlarm,
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

	seriesType := a.seriesTypeFromRequest(r)

	if !a.seriesTypeExists(seriesType) {
		a.notFound(w)

		return
	}

	diskPath := r.URL.Query().Get("path")
	if diskPath == "" {
		a.notFound(w)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeySeriesType, string(seriesType))

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	diskPaths, err := a.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	if !a.diskPathExists(diskPaths, diskPath) {
		a.notFound(w)

		return
	}

	diskMinSeries, diskMaxSeries, diskAvgSeries, err := a.seriesRepository.FindDisk(probe, seriesType, diskPath)
	if err != nil {
		a.serverError(w, err)

		return
	}

	var diskAlarm float64 = 0

	alarm, err := a.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		a.serverError(w, err)

		return
	}

	if alarm != nil {
		diskAlarm = alarm.Disk
	}

	err = a.inertiaManager.Render(w, r, "disk/Index", map[string]any{
		"seriesType":    seriesType,
		"chunkSize":     a.seriesRepository.ChunkSize(seriesType),
		"diskPath":      diskPath,
		"diskPaths":     diskPaths,
		"diskMinSeries": diskMinSeries,
		"diskMaxSeries": diskMaxSeries,
		"diskAvgSeries": diskAvgSeries,
		"diskAlarm":     diskAlarm,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) seriesTypeFromRequest(r *http.Request) models.SeriesType {
	seriesType := models.SeriesType(r.URL.Query().Get("type"))

	if seriesType == "" {
		exists := a.sessionManager.Exists(r.Context(), sessionKeySeriesType)
		if !exists {
			return a.appConfig.SeriesButtons[0]
		}

		return models.SeriesType(a.sessionManager.GetString(r.Context(), sessionKeySeriesType))
	}

	return seriesType
}

func (a *app) diskPathExists(diskPaths []string, diskPath string) bool {
	return slices.Contains(diskPaths, diskPath)
}

func (a *app) seriesTypeExists(seriesType models.SeriesType) bool {
	for _, current := range models.SeriesTypes {
		if current["value"] == seriesType {
			return true
		}
	}

	return false
}
