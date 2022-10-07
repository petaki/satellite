package web

import (
	"log"
	"net/http"

	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/mix"
)

type app struct {
	debug                  bool
	url                    string
	errorLog               *log.Logger
	infoLog                *log.Logger
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
