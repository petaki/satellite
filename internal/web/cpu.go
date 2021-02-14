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

	cpuSeries, err := app.seriesRepository.FindCpu(models.Day)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "cpu/Index", map[string]interface{}{
		"isCpuActive": true,
		"cpuSeries":   cpuSeries,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
