package web

import "net/http"

func (app *App) performanceIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)

		return
	}

	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	err := app.inertiaManager.Render(w, r, "performance/Index", map[string]interface{}{
		"isPerformanceActive": true,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
