package web

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/mix"
)

type app struct {
	debug            bool
	url              string
	errorLog         *log.Logger
	infoLog          *log.Logger
	redisPool        *redis.Pool
	redisKeyPrefix   string
	mixManager       *mix.Mix
	inertiaManager   *inertia.Inertia
	seriesRepository models.SeriesRepository
}
