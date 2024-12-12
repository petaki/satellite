package web

import (
	"context"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/satellite/resources/views"
	"github.com/petaki/satellite/static"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/mix"
)

// Serve function.
func Serve(
	debug bool,
	addr,
	url string,
	seriesButtons []models.SeriesType,
	redisPool *redis.Pool,
	heartbeatEnabled bool,
	heartbeatWait, heartbeatSleep int,
	heartbeatWebhookMethod, heartbeatWebhookURL string,
	heartbeatWebhookHeader map[string]string,
	heartbeatWebhookBody string,
) {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, "satellite:scs:session:")

	mixManager, inertiaManager, err := newMixAndInertiaManager(debug, url)
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	inertiaManager.Share("seriesButtons", seriesButtons)
	inertiaManager.Share("seriesTypes", models.SeriesTypes)

	webApp := &app{
		debug:                  debug,
		url:                    url,
		seriesButtons:          seriesButtons,
		infoLog:                cli.InfoLog,
		errorLog:               cli.ErrorLog,
		sessionManager:         sessionManager,
		heartbeatWait:          heartbeatWait,
		heartbeatSleep:         heartbeatSleep,
		heartbeatWebhookMethod: heartbeatWebhookMethod,
		heartbeatWebhookURL:    heartbeatWebhookURL,
		heartbeatWebhookHeader: heartbeatWebhookHeader,
		heartbeatWebhookBody:   heartbeatWebhookBody,
		mixManager:             mixManager,
		inertiaManager:         inertiaManager,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		probeRepository: &models.RedisProbeRepository{
			RedisPool: redisPool,
		},
		alarmRepository: &models.RedisAlarmRepository{
			RedisPool: redisPool,
		},
		seriesRepository: &models.RedisSeriesRepository{
			RedisPool: redisPool,
		},
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     webApp.errorLog,
		Handler:      webApp.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	var ticker *time.Ticker
	var doneTicker chan bool

	if heartbeatEnabled {
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

	webApp.infoLog.Printf("Starting server on "+cli.Green("%s"), addr)

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
		if heartbeatEnabled {
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

func newMixAndInertiaManager(debug bool, url string) (*mix.Mix, *inertia.Inertia, error) {
	mixManager := mix.New("", "./static", "")

	var version string
	var err error

	if debug {
		version, err = mixManager.Hash("")
		if err != nil {
			return nil, nil, err
		}
	} else {
		version, err = mixManager.HashFromFS("", static.Files)
		if err != nil {
			return nil, nil, err
		}
	}

	inertiaManager := inertia.NewWithFS(url, "app.gohtml", version, views.Templates)
	inertiaManager.Share("title", "Satellite")

	inertiaManager.ShareFunc("asset", func(path string) (string, error) {
		if debug {
			return mixManager.Mix(path, "")
		}

		return "/" + path, nil
	})

	return mixManager, inertiaManager, nil
}
