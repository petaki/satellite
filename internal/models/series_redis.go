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
	seriesCPUKeyPrefix    = "cpu:"
	seriesMemoryKeyPrefix = "memory:"
	seriesDiskKeyPrefix   = "disk:"
)

// RedisSeriesRepository type.
type RedisSeriesRepository struct {
	RedisPool *redis.Pool
}

// FindCPU function.
func (rsr *RedisSeriesRepository) FindCPU(probe Probe, seriesType SeriesType) (Series, error) {
	return rsr.findAllSeries(probe, seriesType, seriesCPUKeyPrefix, "")
}

// FindMemory function.
func (rsr *RedisSeriesRepository) FindMemory(probe Probe, seriesType SeriesType) (Series, error) {
	return rsr.findAllSeries(probe, seriesType, seriesMemoryKeyPrefix, "")
}

// FindDisk function.
func (rsr *RedisSeriesRepository) FindDisk(probe Probe, seriesType SeriesType, path string) (Series, error) {
	return rsr.findAllSeries(probe, seriesType, seriesDiskKeyPrefix, ":"+base64.StdEncoding.EncodeToString([]byte(path)))
}

// FindDiskPaths function.
func (rsr *RedisSeriesRepository) FindDiskPaths(probe Probe) ([]string, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	cursor := 0
	prefix := string(probe) + ":" + seriesDiskKeyPrefix + strconv.FormatInt(today().Unix(), 10) + ":"

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

func (rsr *RedisSeriesRepository) findAllSeries(probe Probe, seriesType SeriesType, prefix, suffix string) (Series, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var series, rawSeries Series

	for _, timestamp := range rsr.timestamps(seriesType) {
		values, err := redis.Strings(
			conn.Do("HGETALL", string(probe)+":"+prefix+strconv.FormatInt(timestamp, 10)+suffix),
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

			rawSeries = append(rawSeries, Value{
				X: x,
				Y: y,
			})
		}
	}

	if len(rawSeries) == 0 {
		return series, nil
	}

	sort.Slice(rawSeries, func(i, j int) bool {
		return rawSeries[i].X > rawSeries[j].X
	})

	var chunks []Series

	chunkSize := rsr.chunkSize(seriesType)

	for chunkSize < len(rawSeries) {
		rawSeries, chunks = rawSeries[chunkSize:], append(chunks, rawSeries[0:chunkSize:chunkSize])
	}

	chunks = append(chunks, rawSeries)

	for _, chunk := range chunks {
		chunkValue := Value{
			X: 0,
			Y: 0,
		}

		for _, value := range chunk {
			chunkValue.X += value.X
			chunkValue.Y += value.Y
		}

		chunkValue.X = chunkValue.X / int64(len(chunk)) * 1000
		chunkValue.Y = chunkValue.Y / float64(len(chunk))

		series = append(series, chunkValue)
	}

	return series, nil
}

func (rsr *RedisSeriesRepository) chunkSize(seriesType SeriesType) int {
	switch seriesType {
	case Week:
		return 90 // 1 hour 30 minutes
	case Month:
		return 60 * 8 // 8 hours
	default:
		return 15 // 15 minutes
	}
}

func (rsr *RedisSeriesRepository) timestamps(seriesType SeriesType) []int64 {
	var timestamps []int64

	end := today()
	var start time.Time

	switch seriesType {
	case Week:
		start = end.AddDate(0, 0, -6)
	case Month:
		start = end.AddDate(0, -1, 0)
	default:
		start = end
	}

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		timestamps = append(timestamps, current.Unix())
	}

	return timestamps
}
