package models

import (
	"github.com/gomodule/redigo/redis"
	"sort"
	"strconv"
	"strings"
)

// RedisProbeRepository type.
type RedisProbeRepository struct {
	RedisPool *redis.Pool
}

// FindAll function.
func (rpr *RedisProbeRepository) FindAll() ([]Probe, error) {
	conn := rpr.RedisPool.Get()
	defer conn.Close()

	cursor := 0
	suffix := ":" + seriesCPUKeyPrefix + strconv.FormatInt(today().Unix(), 10)

	var names []string

	for {
		values, err := redis.Values(
			conn.Do("SCAN", cursor, "MATCH", "*"+suffix),
		)
		if err != nil {
			return nil, err
		}

		cursor, err = redis.Int(values[0], nil)
		if err != nil {
			return nil, err
		}

		current, err := redis.Strings(values[1], nil)
		if err != nil {
			return nil, err
		}

		names = append(names, current...)

		if cursor == 0 {
			break
		}
	}

	sort.Strings(names)

	probes := make([]Probe, len(names))

	for key, value := range names {
		probes[key] = Probe(strings.ReplaceAll(value, suffix, ""))
	}

	return probes, nil
}
