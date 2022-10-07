package web

import (
	"log"

	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/mix"
)

type app struct {
	debug                  bool
	url                    string
	errorLog               *log.Logger
	infoLog                *log.Logger
	heartbeatEnabled       bool
	heartbeatWait          int
	heartbeatSleep         int
	heartbeatWebhookMethod string
	heartbeatWebhookURL    string
	heartbeatWebhookHeader map[string]string
	heartbeatWebhookBody   string
	mixManager             *mix.Mix
	inertiaManager         *inertia.Inertia
	probeRepository        models.ProbeRepository
	alarmRepository        models.AlarmRepository
	seriesRepository       models.SeriesRepository
}
