package models

import (
	"encoding/base64"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	seriesCPUKeyPrefix           = "cpu:"
	seriesMemoryKeyPrefix        = "memory:"
	seriesProcessCPUKeyPrefix    = "process:cpu:"
	seriesProcessMemoryKeyPrefix = "process:memory:"
	seriesLoadKeyPrefix          = "load:"
	seriesDiskKeyPrefix          = "disk:"
	seriesLoadSegmentCount       = 3
	seriesProcessCount           = 3
	seriesProcessSegmentCount    = 2
	seriesScanCount              = 1000
	timestampMultiplier          = 1000
)

type loadSample struct {
	x     int64
	loads [seriesLoadSegmentCount]float64
}

// RedisSeriesRepository type.
type RedisSeriesRepository struct {
	RedisPool *redis.Pool
}

// FindCPU function.
func (rsr *RedisSeriesRepository) FindCPU(probe Probe, seriesType SeriesType) (Series, Series, Series, ProcessSeries, ProcessSeries, ProcessSeries, error) {
	return rsr.findProcessSeries(probe, seriesType, seriesCPUKeyPrefix, seriesProcessCPUKeyPrefix)
}

// FindMemory function.
func (rsr *RedisSeriesRepository) FindMemory(probe Probe, seriesType SeriesType) (Series, Series, Series, ProcessSeries, ProcessSeries, ProcessSeries, error) {
	return rsr.findProcessSeries(probe, seriesType, seriesMemoryKeyPrefix, seriesProcessMemoryKeyPrefix)
}

// FindLoad function.
func (rsr *RedisSeriesRepository) FindLoad(probe Probe, seriesType SeriesType) (Series, Series, Series, error) {
	return rsr.findLoadSeries(probe, seriesType)
}

// FindDisk function.
func (rsr *RedisSeriesRepository) FindDisk(probe Probe, seriesType SeriesType, path string) (Series, Series, Series, error) {
	return rsr.findThresholdSeries(probe, seriesType, seriesDiskKeyPrefix, ":"+base64.StdEncoding.EncodeToString([]byte(path)))
}

// FindDiskPaths function.
func (rsr *RedisSeriesRepository) FindDiskPaths(probe Probe) ([]string, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var paths []string
	timestamps := rsr.timestamps(Last30Days)

	for _, timestamp := range slices.Backward(timestamps) {
		cursor := 0
		prefix := string(probe) + ":" + seriesDiskKeyPrefix + strconv.FormatInt(timestamp, 10) + ":"

		for {
			values, err := redis.Values(
				conn.Do("SCAN", cursor, "MATCH", prefix+"*", "COUNT", seriesScanCount),
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
			path, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
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
	case Last6Hours:
		return 60 // 60 minutes
	case Last12Hours:
		return 60 // 60 minutes
	case Last24Hours:
		return 60 // 60 minutes
	case Last2Days:
		return 60 // 60 minutes
	case Last7Days:
		return 60 * 24 // 24 hours
	case Last30Days:
		return 60 * 24 * 4 // 96 hours
	}

	return 1 // 1 minute
}

// FindLatestCPU function.
func (rsr *RedisSeriesRepository) FindLatestCPU(probe Probe) (float64, bool, error) {
	return rsr.findLatestMetric(probe, seriesCPUKeyPrefix)
}

// FindLatestMemory function.
func (rsr *RedisSeriesRepository) FindLatestMemory(probe Probe) (float64, bool, error) {
	return rsr.findLatestMetric(probe, seriesMemoryKeyPrefix)
}

// FindLatestLoad function.
func (rsr *RedisSeriesRepository) FindLatestLoad(probe Probe) (float64, float64, float64, bool, error) {
	value, found, err := rsr.findLatestValue(probe, seriesLoadKeyPrefix)
	if err != nil || !found {
		return 0, 0, 0, false, err
	}

	loads, ok := parseLoads(value)
	if !ok {
		return 0, 0, 0, false, nil
	}

	return loads[0], loads[1], loads[2], true, nil
}

func (rsr *RedisSeriesRepository) findLatestMetric(probe Probe, prefix string) (float64, bool, error) {
	value, found, err := rsr.findLatestValue(probe, prefix)
	if err != nil || !found {
		return 0, false, err
	}

	y, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false, nil
	}

	return y, true, nil
}

func (rsr *RedisSeriesRepository) findLatestValue(probe Probe, prefix string) (string, bool, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	todayTS := today()
	yesterdayTS := todayTS.AddDate(0, 0, -1)

	for _, ts := range []int64{todayTS.Unix(), yesterdayTS.Unix()} {
		values, err := redis.Strings(
			conn.Do("HGETALL", string(probe)+":"+prefix+strconv.FormatInt(ts, 10)),
		)
		if err != nil {
			return "", false, err
		}

		if len(values) == 0 {
			continue
		}

		var maxTS int64
		var maxValue string

		for i := 0; i < len(values); i += 2 {
			x, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return "", false, err
			}

			if x > maxTS {
				maxTS = x
				maxValue = values[i+1]
			}
		}

		if maxTS == 0 {
			continue
		}

		return maxValue, true, nil
	}

	return "", false, nil
}

func (rsr *RedisSeriesRepository) findLoadSeries(probe Probe, seriesType SeriesType) (Series, Series, Series, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var samples []loadSample

	chunkSize := rsr.ChunkSize(seriesType)
	start := now()
	end := rsr.end(start, seriesType)

	for _, timestamp := range rsr.timestamps(seriesType) {
		values, err := redis.Strings(
			conn.Do("HGETALL", string(probe)+":"+seriesLoadKeyPrefix+strconv.FormatInt(timestamp, 10)),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		for i := 0; i < len(values); i += 2 {
			x, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			if !between(start, end, x) {
				continue
			}

			loads, ok := parseLoads(values[i+1])
			if !ok {
				continue
			}

			samples = append(samples, loadSample{
				x:     x,
				loads: loads,
			})
		}
	}

	if len(samples) == 0 {
		return nil, nil, nil, nil
	}

	sort.SliceStable(samples, func(i, j int) bool {
		return samples[i].x > samples[j].x
	})

	var series [seriesLoadSegmentCount]Series

	for _, chunk := range chunkSlice(chunkSize, samples) {
		var sums [seriesLoadSegmentCount]float64
		var x int64

		for index, sample := range chunk {
			if index == len(chunk)/2 {
				x = sample.x
			}

			for j, load := range sample.loads {
				sums[j] += load
			}
		}

		for j := range series {
			series[j] = append(series[j], Value{
				X: x * timestampMultiplier,
				Y: sums[j] / float64(len(chunk)),
			})
		}
	}

	return series[0], series[1], series[2], nil
}

func (rsr *RedisSeriesRepository) findThresholdSeries(probe Probe, seriesType SeriesType, prefix, suffix string) (Series, Series, Series, error) {
	conn := rsr.RedisPool.Get()
	defer conn.Close()

	var minSeries, maxSeries, avgSeries, rawSeries Series

	chunkSize := rsr.ChunkSize(seriesType)
	start := now()
	end := rsr.end(start, seriesType)

	for _, timestamp := range rsr.timestamps(seriesType) {
		values, err := redis.Strings(
			conn.Do("HGETALL", string(probe)+":"+prefix+strconv.FormatInt(timestamp, 10)+suffix),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		for i := 0; i < len(values); i += 2 {
			x, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			if !between(start, end, x) {
				continue
			}

			y, err := strconv.ParseFloat(values[i+1], 64)
			if err != nil {
				continue
			}

			rawSeries = append(rawSeries, Value{
				X: x,
				Y: y,
			})
		}
	}

	if len(rawSeries) == 0 {
		return nil, nil, nil, nil
	}

	sort.SliceStable(rawSeries, func(i, j int) bool {
		return rawSeries[i].X > rawSeries[j].X
	})

	for _, chunk := range chunkSlice(chunkSize, rawSeries) {
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
		}

		x *= timestampMultiplier

		minValue.X = x
		maxValue.X = x
		avgValue.X = x

		avgValue.Y = avgValue.Y / float64(len(chunk))

		minSeries = append(minSeries, minValue)
		maxSeries = append(maxSeries, maxValue)
		avgSeries = append(avgSeries, avgValue)
	}

	return minSeries, maxSeries, avgSeries, nil
}

func (rsr *RedisSeriesRepository) findProcessSeries(probe Probe, seriesType SeriesType, prefix, processPrefix string) (Series, Series, Series, ProcessSeries, ProcessSeries, ProcessSeries, error) {
	minSeries, maxSeries, avgSeries, err := rsr.findThresholdSeries(probe, seriesType, prefix, "")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	conn := rsr.RedisPool.Get()
	defer conn.Close()

	process1Series := make(ProcessSeries, len(avgSeries))
	process2Series := make(ProcessSeries, len(avgSeries))
	process3Series := make(ProcessSeries, len(avgSeries))

	days := map[int64][]int64{}

	for _, v := range avgSeries {
		timestamp := v.X / timestampMultiplier
		date := day(timestamp).Unix()

		days[date] = append(days[date], timestamp)
	}

	processValues := make(map[int64]string, len(avgSeries))

	for date, timestamps := range days {
		fields := make([]string, len(timestamps))

		for i, timestamp := range timestamps {
			fields[i] = strconv.FormatInt(timestamp, 10)
		}

		values, err := redis.Values(
			conn.Do("HMGET", redis.Args{}.Add(string(probe)+":"+processPrefix+strconv.FormatInt(date, 10)).AddFlat(fields)...),
		)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}

		for i, value := range values {
			if i >= len(timestamps) {
				break
			}

			raw, err := redis.String(value, nil)
			if err != nil {
				continue
			}

			processValues[timestamps[i]] = raw
		}
	}

	placeholder := now().UnixMilli()

	for k, v := range avgSeries {
		column := make(ProcessSeries, seriesProcessCount)

		for i := range column {
			column[i] = ProcessValue{
				Name: "Not Set",
				X:    placeholder,
				Y:    0,
			}
		}

		for j, process := range strings.Split(processValues[v.X/timestampMultiplier], "|") {
			if j >= len(column) {
				break
			}

			segments := strings.SplitN(process, ":", seriesProcessSegmentCount)

			if len(segments) != seriesProcessSegmentCount {
				continue
			}

			py, err := strconv.ParseFloat(segments[1], 64)
			if err != nil {
				continue
			}

			column[j] = ProcessValue{
				Name: segments[0],
				X:    v.X,
				Y:    py,
			}
		}

		process1Series[k] = column[0]
		process2Series[k] = column[1]
		process3Series[k] = column[2]
	}

	sort.SliceStable(process1Series, func(i, j int) bool {
		return process1Series[i].X > process1Series[j].X
	})

	sort.SliceStable(process2Series, func(i, j int) bool {
		return process2Series[i].X > process2Series[j].X
	})

	sort.SliceStable(process3Series, func(i, j int) bool {
		return process3Series[i].X > process3Series[j].X
	})

	return minSeries, maxSeries, avgSeries, process1Series, process2Series, process3Series, nil
}

func parseLoads(value string) ([seriesLoadSegmentCount]float64, bool) {
	var loads [seriesLoadSegmentCount]float64

	segments := strings.SplitN(value, ":", seriesLoadSegmentCount)

	if len(segments) != seriesLoadSegmentCount {
		return loads, false
	}

	for i, segment := range segments {
		load, err := strconv.ParseFloat(segment, 64)
		if err != nil {
			return loads, false
		}

		loads[i] = load
	}

	return loads, true
}

func chunkSlice[T any](chunkSize int, values []T) [][]T {
	var chunks [][]T

	for chunkSize < len(values) {
		values, chunks = values[chunkSize:], append(chunks, values[0:chunkSize:chunkSize])
	}

	return append(chunks, values)
}

func (rsr *RedisSeriesRepository) end(start time.Time, seriesType SeriesType) time.Time {
	switch seriesType {
	case Last15Minutes:
		return start.Add(-15 * time.Minute)
	case Last30Minutes:
		return start.Add(-30 * time.Minute)
	case Last1Hour:
		return start.Add(-1 * time.Hour)
	case Last3Hours:
		return start.Add(-3 * time.Hour)
	case Last6Hours:
		return start.Add(-6 * time.Hour)
	case Last12Hours:
		return start.Add(-12 * time.Hour)
	case Last24Hours:
		return start.Add(-24 * time.Hour)
	case Last2Days:
		return start.Add(-2 * 24 * time.Hour)
	case Last7Days:
		return start.Add(-7 * 24 * time.Hour)
	case Last30Days:
		return start.Add(-30 * 24 * time.Hour)
	}

	return start.Add(-5 * time.Minute)
}

func (rsr *RedisSeriesRepository) timestamps(seriesType SeriesType) []int64 {
	var timestamps []int64

	end := today()
	var start time.Time

	switch seriesType {
	case Last2Days:
		start = end.AddDate(0, 0, -2)
	case Last7Days:
		start = end.AddDate(0, 0, -7)
	case Last30Days:
		start = end.AddDate(0, 0, -30)
	default:
		start = end.AddDate(0, 0, -1)
	}

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		timestamps = append(timestamps, current.Unix())
	}

	return timestamps
}
