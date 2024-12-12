package web

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/mix"
)

const (
	sessionKeySeriesType = "seriesType"
)

type app struct {
	debug                  bool
	url                    string
	seriesButtons          []models.SeriesType
	infoLog                *log.Logger
	errorLog               *log.Logger
	sessionManager         *scs.SessionManager
	heartbeatWait          int
	heartbeatSleep         int
	heartbeatWebhookMethod string
	heartbeatWebhookURL    string
	heartbeatWebhookHeader map[string]string
	heartbeatWebhookBody   string
	mixManager             *mix.Mix
	inertiaManager         *inertia.Inertia
	client                 *http.Client
	probeRepository        models.ProbeRepository
	alarmRepository        models.AlarmRepository
	seriesRepository       models.SeriesRepository
}
