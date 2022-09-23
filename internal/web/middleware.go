package web

import (
	"context"
	"fmt"

	"github.com/petaki/satellite/internal/models"
	"golang.org/x/exp/slices"
	"net/http"
)

func (a *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *app) probes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		probes, err := a.probeRepository.FindAll()
		if err != nil {
			a.serverError(w, err)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyProbes, probes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *app) probe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("probe") == "" {
			a.notFound(w)

			return
		}

		probe := models.Probe(r.URL.Query().Get("probe"))

		if !slices.Contains(r.Context().Value(contextKeyProbes).([]models.Probe), probe) {
			a.notFound(w)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyProbe, probe)
		ctx = a.inertiaManager.WithProp(ctx, "probe", probe)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
