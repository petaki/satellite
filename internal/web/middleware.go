package web

import (
	"context"
	"net/http"
	"slices"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) probe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("probe") == "" {
			a.notFound(w)

			return
		}

		probe := models.Probe(r.URL.Query().Get("probe"))

		probes, err := a.probeRepository.FindAll()
		if err != nil {
			a.serverError(w, err)

			return
		}

		if !slices.Contains(probes, probe) {
			a.notFound(w)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyProbe, probe)
		ctx = a.inertiaManager.WithProp(ctx, "probe", probe)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *app) allowedHosts(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostname := requestHostname(r.Host)

		if !slices.Contains(a.appConfig.AllowedHosts, "*") &&
			!slices.Contains(a.appConfig.AllowedHosts, hostname) &&
			!isLoopback(hostname) {
			a.forbidden(w)

			return
		}

		next.ServeHTTP(w, r)
	})
}
