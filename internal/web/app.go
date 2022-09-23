package web

import (
	"log"

	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/mix"
)

type app struct {
	debug            bool
	url              string
	errorLog         *log.Logger
	infoLog          *log.Logger
	mixManager       *mix.Mix
	inertiaManager   *inertia.Inertia
	probeRepository  models.ProbeRepository
	alarmRepository  models.AlarmRepository
	seriesRepository models.SeriesRepository
}
