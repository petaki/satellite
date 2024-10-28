package models

import (
	"cmp"
	"encoding/base64"
	"golang.org/x/exp/slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	seriesCPUKeyPrefix           = "cpu:"
	seriesMemoryKeyPrefix        = "memory:"
	seriesDiskKeyPrefix          = "disk:"
	seriesProcessCPUKeyPrefix    = "process:cpu:"
	seriesProcessMemoryKeyPrefix = "process:memory:"
)

// RedisSeriesRepository type.
type RedisSeriesRepository struct {
	RedisPool *redis.Pool
}

// FindCPU function.
func (rsr *RedisSeriesRepository) FindCPU(probe Probe, seriesType SeriesType) (Series, Series, Series, error) {
	return rsr.findAllSeries(probe, seriesType, CPU, seriesCPUKeyPrefix, "")
}

// FindMemory function.
func (rsr *RedisSeriesRepository) FindMemory(probe Probe, seriesType SeriesType) (Series, Series, Series, error) {
	return rsr.findAllSeries(probe, seriesType, Memory, seriesMemoryKeyPrefix, "")
}

// FindDisk function.
func (rsr *RedisSeriesRepository) FindDisk(probe Probe, seriesType SeriesType, path string) (Series, Series, Series, error) {
	return rsr.findAllSeries(probe, seriesType, Disk, seriesDiskKeyPrefix, ":"+base64.StdEncoding.EncodeToString([]byte(path)))
}

// FindDiskPaths function.
func (rsr *RedisSeriesRepository) FindDiskPaths(probe Probe) ([]string, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var paths []string
	timestamps := rsr.timestamps(Month)

	for i := len(timestamps) - 1; i >= 0; i-- {
		cursor := 0
		prefix := string(probe) + ":" + seriesDiskKeyPrefix + strconv.FormatInt(timestamps[i], 10) + ":"

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

		if len(paths) == 0 {
			continue
		}

		for key, value := range paths {
			path, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(value, prefix, ""))
			if err != nil {
				return nil, err
			}

			paths[key] = string(path)
		}

		sort.Strings(paths)

		break
	}

	return paths, nil
}

// ChunkSize function.
func (rsr *RedisSeriesRepository) ChunkSize(seriesType SeriesType) int {
	switch seriesType {
	case Week:
		return 60 * 24 // 24 hours
	case Month:
		return 60 * 24 // 24 hours
	}

	return 1 // 1 minute
}

func (rsr *RedisSeriesRepository) findAllSeries(probe Probe, seriesType SeriesType, watcherType WatcherType, prefix, suffix string) (Series, Series, Series, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var minSeries, maxSeries, avgSeries, rawSeries Series

	chunkSize := rsr.ChunkSize(seriesType)

	for _, timestamp := range rsr.timestamps(seriesType) {
		rawSeries = nil

		values, err := redis.Strings(
			conn.Do("HGETALL", string(probe)+":"+prefix+strconv.FormatInt(timestamp, 10)+suffix),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		var processValues []string

		if watcherType == CPU {
			processValues, err = redis.Strings(
				conn.Do("HGETALL", string(probe)+":"+seriesProcessCPUKeyPrefix+strconv.FormatInt(timestamp, 10)),
			)
			if err != nil {
				processValues = nil
			}
		}

		if watcherType == Memory {
			processValues, err = redis.Strings(
				conn.Do("HGETALL", string(probe)+":"+seriesProcessMemoryKeyPrefix+strconv.FormatInt(timestamp, 10)),
			)
			if err != nil {
				processValues = nil
			}
		}

		for i := 0; i < len(values); i += 2 {
			x, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			y, err := strconv.ParseFloat(values[i+1], 64)
			if err != nil {
				return nil, nil, nil, err
			}

			v := Value{
				X: x,
				Y: y,
			}

			if (watcherType == CPU || watcherType == Memory) && len(processValues) > i+1 {
				processes := strings.Split(processValues[i+1], "|")

				for _, p := range processes {
					segments := strings.SplitN(p, ":", 2)

					py, err := strconv.ParseFloat(segments[1], 64)
					if err != nil {
						continue
					}

					v.P = append(v.P, ProcessValue{
						Name: segments[0],
						Y:    py,
					})
				}
			}

			rawSeries = append(rawSeries, v)
		}

		if len(rawSeries) == 0 {
			continue
		}

		sort.SliceStable(rawSeries, func(i, j int) bool {
			return rawSeries[i].X > rawSeries[j].X
		})

		for _, chunk := range rsr.chunks(chunkSize, rawSeries) {
			minValue := Value{
				X: 0,
				Y: 0,
			}

			maxValue := Value{
				X: 0,
				Y: 0,
			}

			avgValue := Value{
				X: 0,
				Y: 0,
			}

			var x int64 = 0

			for index, value := range chunk {
				if index == len(chunk)/2 {
					x = value.X
				}

				if index == 0 {
					minValue.Y = value.Y
					maxValue.Y = value.Y
				} else {
					if minValue.Y > value.Y {
						minValue.Y = value.Y
					}

					if maxValue.Y < value.Y {
						maxValue.Y = value.Y
					}
				}

				avgValue.Y += value.Y

				if watcherType == CPU || watcherType == Memory {
					for _, process := range value.P {
						index := slices.IndexFunc(avgValue.P, func(pv ProcessValue) bool {
							return process.Name == pv.Name
						})

						if index == -1 {
							avgValue.P = append(avgValue.P, process)
						} else {
							avgValue.P[index].Y = (avgValue.P[index].Y + process.Y) / 2
						}
					}
				}
			}

			x *= 1000

			minValue.X = x
			maxValue.X = x
			avgValue.X = x

			avgValue.Y = avgValue.Y / float64(len(chunk))

			if watcherType == CPU || watcherType == Memory {
				slices.SortStableFunc(avgValue.P, func(a, b ProcessValue) int {
					return cmp.Compare(b.Y, a.Y)
				})
			}

			minSeries = append(minSeries, minValue)
			maxSeries = append(maxSeries, maxValue)
			avgSeries = append(avgSeries, avgValue)
		}
	}

	return minSeries, maxSeries, avgSeries, nil
}

func (rsr *RedisSeriesRepository) chunks(chunkSize int, series Series) []Series {
	var chunks []Series

	for chunkSize < len(series) {
		series, chunks = series[chunkSize:], append(chunks, series[0:chunkSize:chunkSize])
	}

	return append(chunks, series)
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
