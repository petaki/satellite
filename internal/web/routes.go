package web

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/petaki/satellite/static"
)

func (a *app) routes() http.Handler {
	baseMiddleware := alice.New(a.recoverPanic)
	webMiddleware := alice.New(
		a.sessionManager.LoadAndSave,
		a.probes,
		a.inertiaManager.Middleware,
	)

	mux := http.NewServeMux()
	mux.Handle("/", webMiddleware.ThenFunc(a.probeIndex))
	mux.Handle("/cpu", webMiddleware.Append(a.probe).ThenFunc(a.cpuIndex))
	mux.Handle("/memory", webMiddleware.Append(a.probe).ThenFunc(a.memoryIndex))
	mux.Handle("/load", webMiddleware.Append(a.probe).ThenFunc(a.loadIndex))
	mux.Handle("/disk", webMiddleware.Append(a.probe).ThenFunc(a.diskIndex))
	mux.Handle("/probe/delete", webMiddleware.Append(a.probe).ThenFunc(a.probeDelete))
	mux.Handle("/probe/delete-all", webMiddleware.ThenFunc(a.probeDeleteAll))

	var fileServer http.Handler

	if a.debug {
		fileServer = http.FileServer(http.Dir("./static/"))
	} else {
		staticFS := http.FS(static.Files)
		fileServer = http.FileServer(staticFS)
	}

	mux.Handle("/css/", fileServer)
	mux.Handle("/images/", fileServer)
	mux.Handle("/js/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	return baseMiddleware.Then(mux)
}
