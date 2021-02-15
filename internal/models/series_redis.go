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
	return rsr.findAllSeries(seriesType, seriesCpuKeyPrefix, "")
}

func (rsr *RedisSeriesRepository) FindMemory(seriesType SeriesType) (Series, error) {
	return rsr.findAllSeries(seriesType, seriesMemoryKeyPrefix, "")
}

func (rsr *RedisSeriesRepository) FindDisk(seriesType SeriesType, path string) (Series, error) {
	return rsr.findAllSeries(seriesType, seriesDiskKeyPrefix, ":"+base64.StdEncoding.EncodeToString([]byte(path)))
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

func (rsr *RedisSeriesRepository) findAllSeries(seriesType SeriesType, prefix, suffix string) (Series, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var series Series

	for _, timestamp := range rsr.timestamps(seriesType) {
		values, err := redis.Strings(
			conn.Do("HGETALL", rsr.RedisKeyPrefix+prefix+strconv.FormatInt(timestamp, 10)+suffix),
		)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(values); i += 2 {
			x, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return nil, err
			}

			y, err := strconv.ParseFloat(values[i+1], 64)
			if err != nil {
				return nil, err
			}

			series = append(series, Value{
				X: x,
				Y: y,
			})
		}
	}

	return series, nil
}

func (rsr *RedisSeriesRepository) timestamps(seriesType SeriesType) []int64 {
	var timestamps []int64

	end := rsr.today()
	var start time.Time

	switch seriesType {
	case Week:
		start = end.AddDate(0, 0, -6)
	case Month:
		start = end.AddDate(0, -1, 0)
	default:
		start = end
	}

	for current := start; current.After(end) == false; current = current.AddDate(0, 0, 1) {
		timestamps = append(timestamps, current.Unix())
	}

	return timestamps
}

func (rsr *RedisSeriesRepository) today() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}
