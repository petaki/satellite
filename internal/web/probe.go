package web

import (
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

	err := a.inertiaManager.Render(w, r, "probe/Index", map[string]interface{}{
		"isProbeActive": true,
		"probes":        r.Context().Value(contextKeyProbes).([]models.Probe),
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
