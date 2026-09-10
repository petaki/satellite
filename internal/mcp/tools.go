package mcp

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
)

type handler struct {
	appConfig        *config.Config
	probeRepository  models.ProbeRepository
	alarmRepository  models.AlarmRepository
	seriesRepository models.SeriesRepository
	logRepository    models.LogRepository
}

func (h *handler) listProbes(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probes, err := h.probeRepository.FindAll()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: probes: %v", ErrFind, err)), nil
	}

	summaries := make([]models.ProbeSummary, 0, len(probes))

	for _, probe := range probes {
		summary := models.ProbeSummary{
			Name: string(probe),
		}

		cpu, cpuFound, err := h.seriesRepository.FindLatestCPU(probe)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: cpu, %s: %v", ErrFind, probe, err)), nil
		}

		if cpuFound {
			summary.CPU = cpu
		}

		mem, memFound, err := h.seriesRepository.FindLatestMemory(probe)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: memory, %s: %v", ErrFind, probe, err)), nil
		}

		if memFound {
			summary.Memory = mem
		}

		load1, load5, load15, _, err := h.seriesRepository.FindLatestLoad(probe)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: load, %s: %v", ErrFind, probe, err)), nil
		}

		summary.Load1 = load1
		summary.Load5 = load5
		summary.Load15 = load15

		values, _, err := h.probeRepository.FindLatestValues(probe, 2)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: latest values, %s: %v", ErrFind, probe, err)), nil
		}

		for _, value := range values {
			if value != nil {
				summary.IsActive = true

				break
			}
		}

		alarm, err := h.alarmRepository.Find(probe)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
		}

		if alarm != nil {
			summary.CPUAlarm = alarm.CPU
			summary.MemAlarm = alarm.Memory
			summary.LoadAlarm = alarm.Load
		}

		summaries = append(summaries, summary)
	}

	return marshalResult(summaries)
}

func (h *handler) getCPU(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	seriesType := h.resolveSeriesType(request)
	if !h.seriesTypeExists(seriesType) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: series type, %s", ErrInvalid, seriesType)), nil
	}

	cpuMin, cpuMax, cpuAvg, proc1, proc2, proc3, err := h.seriesRepository.FindCPU(probe, seriesType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: cpu, %s, %s: %v", ErrFind, probe, seriesType, err)), nil
	}

	var cpuAlarm float64

	alarm, err := h.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
	}

	if alarm != nil {
		cpuAlarm = alarm.CPU
	}

	result := map[string]any{
		"probe":          string(probe),
		"series":         string(seriesType),
		"cpuMinSeries":   cpuMin,
		"cpuMaxSeries":   cpuMax,
		"cpuAvgSeries":   cpuAvg,
		"process1Series": proc1,
		"process2Series": proc2,
		"process3Series": proc3,
		"cpuAlarm":       cpuAlarm,
	}

	return marshalResult(result)
}

func (h *handler) getMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	seriesType := h.resolveSeriesType(request)
	if !h.seriesTypeExists(seriesType) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: series type, %s", ErrInvalid, seriesType)), nil
	}

	memMin, memMax, memAvg, proc1, proc2, proc3, err := h.seriesRepository.FindMemory(probe, seriesType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: memory, %s, %s: %v", ErrFind, probe, seriesType, err)), nil
	}

	var memoryAlarm float64

	alarm, err := h.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
	}

	if alarm != nil {
		memoryAlarm = alarm.Memory
	}

	result := map[string]any{
		"probe":           string(probe),
		"series":          string(seriesType),
		"memoryMinSeries": memMin,
		"memoryMaxSeries": memMax,
		"memoryAvgSeries": memAvg,
		"process1Series":  proc1,
		"process2Series":  proc2,
		"process3Series":  proc3,
		"memoryAlarm":     memoryAlarm,
	}

	return marshalResult(result)
}

func (h *handler) getLoad(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	seriesType := h.resolveSeriesType(request)
	if !h.seriesTypeExists(seriesType) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: series type, %s", ErrInvalid, seriesType)), nil
	}

	load1, load5, load15, err := h.seriesRepository.FindLoad(probe, seriesType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: load, %s, %s: %v", ErrFind, probe, seriesType, err)), nil
	}

	var loadAlarm float64

	alarm, err := h.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
	}

	if alarm != nil {
		loadAlarm = alarm.Load
	}

	result := map[string]any{
		"probe":        string(probe),
		"series":       string(seriesType),
		"load1Series":  load1,
		"load5Series":  load5,
		"load15Series": load15,
		"loadAlarm":    loadAlarm,
	}

	return marshalResult(result)
}

func (h *handler) getDisk(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	seriesType := h.resolveSeriesType(request)
	if !h.seriesTypeExists(seriesType) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: series type, %s", ErrInvalid, seriesType)), nil
	}

	diskPaths, err := h.seriesRepository.FindDiskPaths(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: disk paths, %s: %v", ErrFind, probe, err)), nil
	}

	diskPath := request.GetString("path", "")

	if diskPath != "" && !slices.Contains(diskPaths, diskPath) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: disk path, %s, available: %v", ErrInvalid, diskPath, diskPaths)), nil
	}

	if diskPath == "" && len(diskPaths) > 0 {
		diskPath = diskPaths[0]
	}

	diskMin, diskMax, diskAvg := models.Series{}, models.Series{}, models.Series{}

	if diskPath != "" {
		diskMin, diskMax, diskAvg, err = h.seriesRepository.FindDisk(probe, seriesType, diskPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: disk, %s, %s: %v", ErrFind, probe, diskPath, err)), nil
		}
	}

	var diskAlarm float64

	alarm, err := h.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
	}

	if alarm != nil {
		diskAlarm = alarm.Disk
	}

	result := map[string]any{
		"probe":         string(probe),
		"series":        string(seriesType),
		"diskPath":      diskPath,
		"diskPaths":     diskPaths,
		"diskMinSeries": diskMin,
		"diskMaxSeries": diskMax,
		"diskAvgSeries": diskAvg,
		"diskAlarm":     diskAlarm,
	}

	return marshalResult(result)
}

func (h *handler) getLogs(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	logPaths, err := h.logRepository.FindLogPaths(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: log paths, %s: %v", ErrFind, probe, err)), nil
	}

	logPath := request.GetString("path", "")

	if logPath != "" && !slices.Contains(logPaths, logPath) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: log path, %s, available: %v", ErrInvalid, logPath, logPaths)), nil
	}

	if logPath == "" && len(logPaths) > 0 {
		logPath = logPaths[0]
	}

	logEntries := []models.LogEntry{}

	if logPath != "" {
		logEntries, err = h.logRepository.FindLog(probe, logPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v: log, %s, %s: %v", ErrFind, probe, logPath, err)), nil
		}
	}

	result := map[string]any{
		"probe":      string(probe),
		"logPath":    logPath,
		"logPaths":   logPaths,
		"logEntries": logEntries,
	}

	return marshalResult(result)
}

func (h *handler) getAlerts(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	alarm, err := h.alarmRepository.Find(probe)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return mcp.NewToolResultError(fmt.Sprintf("%v: alarm, %s: %v", ErrFind, probe, err)), nil
	}

	alarms := map[string]float64{}

	if alarm != nil {
		alarms = map[string]float64{
			"cpu":    alarm.CPU,
			"memory": alarm.Memory,
			"disk":   alarm.Disk,
			"load":   alarm.Load,
		}
	}

	cpu, cpuFound, err := h.seriesRepository.FindLatestCPU(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: cpu, %s: %v", ErrFind, probe, err)), nil
	}

	mem, memFound, err := h.seriesRepository.FindLatestMemory(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: memory, %s: %v", ErrFind, probe, err)), nil
	}

	load1, load5, load15, loadFound, err := h.seriesRepository.FindLatestLoad(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: load, %s: %v", ErrFind, probe, err)), nil
	}

	latest := map[string]any{}

	if cpuFound {
		latest["cpu"] = cpu
	}

	if memFound {
		latest["memory"] = mem
	}

	if loadFound {
		latest["load1"] = load1
		latest["load5"] = load5
		latest["load15"] = load15
	}

	result := map[string]any{
		"probe":  string(probe),
		"alarms": alarms,
		"latest": latest,
	}

	return marshalResult(result)
}

func (h *handler) deleteProbe(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	probe, err := h.requireProbe(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = h.probeRepository.Delete(probe)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%v: probe, %s: %v", ErrDelete, probe, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("probe %q deleted successfully", string(probe))), nil
}

func (h *handler) requireProbe(request mcp.CallToolRequest) (models.Probe, error) {
	probeName := request.GetString("probe", "")
	if probeName == "" {
		return "", fmt.Errorf("%w: probe", ErrRequired)
	}

	probes, err := h.probeRepository.FindAll()
	if err != nil {
		return "", fmt.Errorf("%w: probes: %v", ErrFind, err)
	}

	probe := models.Probe(probeName)

	if !slices.Contains(probes, probe) {
		return "", fmt.Errorf("%w: probe, %s", ErrNotFound, probeName)
	}

	return probe, nil
}

func (h *handler) resolveSeriesType(request mcp.CallToolRequest) models.SeriesType {
	series := request.GetString("series", "")
	if series == "" {
		return h.appConfig.SeriesButtons[0]
	}

	return models.SeriesType(series)
}

func (h *handler) seriesTypeExists(seriesType models.SeriesType) bool {
	for _, current := range models.SeriesTypes {
		if current["value"] == seriesType {
			return true
		}
	}

	return false
}
