package web

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	baseMiddleware := alice.New(app.recoverPanic)
	webMiddleware := alice.New(
		app.inertiaManager.Middleware,
	)

	mux := http.NewServeMux()
	mux.Handle("/", webMiddleware.ThenFunc(app.cpuIndex))

	fileServer := http.FileServer(http.Dir("./public/"))

	mux.Handle("/css/", fileServer)
	mux.Handle("/images/", fileServer)
	mux.Handle("/js/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	return baseMiddleware.Then(mux)
}
