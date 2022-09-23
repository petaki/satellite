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
