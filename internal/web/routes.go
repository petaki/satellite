package web

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/petaki/satellite/static"
	supportweb "github.com/petaki/support-go/web"
)

func (a *app) routes() http.Handler {
	baseMiddleware := alice.New(supportweb.RecoverPanic(func(w http.ResponseWriter, r *http.Request, err error) {
		a.serverError(w, err)
	}))
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
	mux.Handle("/log", webMiddleware.Append(a.probe).ThenFunc(a.logIndex))
	mux.Handle("/probe/delete", webMiddleware.Append(a.probe).ThenFunc(a.probeDelete))
	mux.Handle("/probe/delete-all", webMiddleware.ThenFunc(a.probeDeleteAll))

	fileServer := supportweb.Static(a.appConfig.Debug, "./static/", static.Files)

	mux.Handle("/build/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	if a.mcpServer != nil {
		mux.Handle("/mcp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				rc := http.NewResponseController(w)
				_ = rc.SetWriteDeadline(time.Time{})
			}

			a.mcpServer.ServeHTTP(w, r)
		}))
	}

	return baseMiddleware.Then(mux)
}
