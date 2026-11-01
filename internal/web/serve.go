package web

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/config"
	satellitemcp "github.com/petaki/satellite/internal/mcp"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/satellite/internal/service"
	"github.com/petaki/satellite/resources/views"
	"github.com/petaki/satellite/static"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/vite"
	supportweb "github.com/petaki/support-go/web"
)

// Serve function.
func Serve(cliApp *cli.App, appConfig *config.Config) error {
	redisPool := service.RedisPool(appConfig)
	defer redisPool.Close()

	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, "satellite:scs:session:")

	viteManager, inertiaManager, err := newViteAndInertiaManager(appConfig)
	if err != nil {
		return err
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

	if appConfig.MCPEnabled {
		webApp.mcpServer = satellitemcp.NewHandler(
			cliApp,
			appConfig,
			webApp.probeRepository,
			webApp.alarmRepository,
			webApp.seriesRepository,
			webApp.logRepository,
		)
	}

	srv := supportweb.NewServer(appConfig.Addr, webApp.routes(), webApp.errorLog)

	if webApp.mcpServer != nil {
		srv.RegisterOnShutdown(func() {
			webApp.mcpServer.CloseSessions(context.Background())
		})
	}

	var heartbeat sync.WaitGroup
	defer heartbeat.Wait()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if appConfig.HeartbeatEnabled {
		heartbeat.Go(func() {
			webApp.runHeartbeat(ctx)
		})
	}

	return srv.RunAndShutdown(ctx, 5*time.Second)
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

	for name, helper := range viteManager.Helpers() {
		inertiaManager.ShareFunc(name, helper)
	}

	return viteManager, inertiaManager, nil
}
