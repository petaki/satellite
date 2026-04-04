package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/satellite/internal/service"
	"github.com/petaki/satellite/resources/views"
	"github.com/petaki/satellite/static"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/vite"
)

// Serve function.
func Serve(cliApp *cli.App, appConfig *config.Config) {
	redisPool := service.RedisPool(appConfig)
	defer redisPool.Close()

	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, "satellite:scs:session:")

	viteManager, inertiaManager, err := newViteAndInertiaManager(appConfig)
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	webApp := &app{
		cliApp:         cliApp,
		appConfig:      appConfig,
		infoLog:        cli.InfoLog,
		errorLog:       cli.ErrorLog,
		sessionManager: sessionManager,
		viteManager:    viteManager,
		inertiaManager: inertiaManager,
		client:         service.HTTPClient(),
		probeRepository: &models.RedisProbeRepository{
			RedisPool: redisPool,
		},
		alarmRepository: &models.RedisAlarmRepository{
			RedisPool: redisPool,
		},
		seriesRepository: &models.RedisSeriesRepository{
			RedisPool: redisPool,
		},
		logRepository: &models.RedisLogRepository{
			RedisPool: redisPool,
		},
	}

	srv := &http.Server{
		Addr:         appConfig.Addr,
		ErrorLog:     webApp.errorLog,
		Handler:      webApp.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	var ticker *time.Ticker
	var doneTicker chan bool

	if appConfig.HeartbeatEnabled {
		ticker = time.NewTicker(time.Minute)
		doneTicker = make(chan bool)

		go func() {
			for {
				select {
				case <-doneTicker:
					return
				case <-ticker.C:
					webApp.heartbeat()
				}
			}
		}()
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	webApp.infoLog.Printf("Starting server on "+cli.Green("%s"), appConfig.Addr)

	go func() {
		err = srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			webApp.errorLog.Fatal(err)
		}
	}()

	<-done
	webApp.infoLog.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if appConfig.HeartbeatEnabled {
			ticker.Stop()
			doneTicker <- true
			webApp.infoLog.Print("Ticker stopped")
		}

		cancel()
	}()

	err = srv.Shutdown(ctx)
	if err != nil {
		webApp.errorLog.Fatal(err)
	}

	webApp.infoLog.Print("Server exited properly")
}

func newViteAndInertiaManager(appConfig *config.Config) (*vite.Vite, *inertia.Inertia, error) {
	var viteManager *vite.Vite
	var version string
	var err error

	if appConfig.Debug {
		viteManager = vite.New("static", "build")
	} else {
		viteManager = vite.New("static", "build", static.Files)
	}

	version, err = viteManager.ManifestHash()
	if err != nil {
		return nil, nil, err
	}

	inertiaManager := inertia.New(appConfig.URL, "app.gohtml", version, views.Templates)
	inertiaManager.Share("title", "Satellite")

	suffix := ""

	if appConfig.Name != "" {
		suffix = fmt.Sprintf(": %s", appConfig.Name)
	}

	inertiaManager.Share("suffix", suffix)
	inertiaManager.Share("seriesButtons", appConfig.SeriesButtons)
	inertiaManager.Share("seriesTypes", models.SeriesTypes)
	inertiaManager.ShareFunc("isRunningHot", viteManager.IsRunningHot)
	inertiaManager.ShareFunc("asset", viteManager.Asset)
	inertiaManager.ShareFunc("css", viteManager.CSS)

	return viteManager, inertiaManager, nil
}
