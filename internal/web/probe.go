package web

import (
	"errors"
	"net/http"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) probeIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)

		return
	}

	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	probes := r.Context().Value(contextKeyProbes).([]models.Probe)

	summaries := make([]models.ProbeSummary, 0, len(probes))

	for _, probe := range probes {
		summary := models.ProbeSummary{
			Name: string(probe),
		}

		cpu, cpuFound, err := a.seriesRepository.FindLatestCPU(probe)
		if err != nil {
			a.serverError(w, err)

			return
		}

		if cpuFound {
			summary.CPU = cpu
		}

		mem, memFound, err := a.seriesRepository.FindLatestMemory(probe)
		if err != nil {
			a.serverError(w, err)

			return
		}

		if memFound {
			summary.Memory = mem
		}

		load1, load5, load15, _, err := a.seriesRepository.FindLatestLoad(probe)
		if err != nil {
			a.serverError(w, err)

			return
		}

		summary.Load1 = load1
		summary.Load5 = load5
		summary.Load15 = load15

		values, _, err := a.probeRepository.FindLatestValues(probe, 2)
		if err != nil {
			a.serverError(w, err)

			return
		}

		for _, value := range values {
			if value != nil {
				summary.IsActive = true

				break
			}
		}

		alarm, err := a.alarmRepository.Find(probe)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			a.serverError(w, err)

			return
		}

		if alarm != nil {
			summary.CPUAlarm = alarm.CPU
			summary.MemAlarm = alarm.Memory
			summary.LoadAlarm = alarm.Load
		}

		summaries = append(summaries, summary)
	}

	err := a.inertiaManager.Render(w, r, "probe/Index", map[string]any{
		"isProbeActive": true,
		"probes":        summaries,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) probeDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		a.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	probe := r.Context().Value(contextKeyProbe).(models.Probe)

	err := a.probeRepository.Delete(probe)
	if err != nil {
		a.serverError(w, err)

		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *app) probeDeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		a.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	probes, err := a.probeRepository.FindAll()
	if err != nil {
		a.serverError(w, err)

		return
	}

	for _, probe := range probes {
		err := a.probeRepository.Delete(probe)
		if err != nil {
			a.serverError(w, err)

			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
