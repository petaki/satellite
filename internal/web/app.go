package web

import (
	"github.com/petaki/satellite/internal/models"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

type App struct {
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
