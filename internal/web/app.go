package web

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/vite"
)

const (
	sessionKeySeriesType = "seriesType"
)

type app struct {
	cliApp           *cli.App
	appConfig        *config.Config
	infoLog          *log.Logger
	errorLog         *log.Logger
	sessionManager   *scs.SessionManager
	viteManager      *vite.Vite
	inertiaManager   *inertia.Inertia
	client           *http.Client
	probeRepository  models.ProbeRepository
	alarmRepository  models.AlarmRepository
	seriesRepository models.SeriesRepository
	logRepository    models.LogRepository
}
