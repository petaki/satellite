package models

import (
	"encoding/base64"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	seriesCpuKeyPrefix    = "cpu:"
	seriesMemoryKeyPrefix = "memory:"
	seriesDiskKeyPrefix   = "disk:"
)

type RedisSeriesRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

func (rsr *RedisSeriesRepository) FindCpu(seriesType SeriesType) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindMemory(seriesType SeriesType) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindDisk(seriesType SeriesType, path string) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindDiskPaths() ([]string, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	cursor := 0
	prefix := rsr.RedisKeyPrefix + seriesDiskKeyPrefix + strconv.FormatInt(rsr.today().Unix(), 10) + ":"

	var paths []string

	for {
		values, err := redis.Values(
			conn.Do("SCAN", cursor, "MATCH", prefix+"*"),
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

		paths = append(paths, current...)

		if cursor == 0 {
			break
		}
	}

	for key, value := range paths {
		path, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(value, prefix, ""))
		if err != nil {
			return nil, err
		}

		paths[key] = string(path)
	}

	sort.Strings(paths)

	return paths, nil
}

func (rsr *RedisSeriesRepository) today() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}
