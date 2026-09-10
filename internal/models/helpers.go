package models

import (
	"encoding/base64"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

func between(start, end time.Time, timestamp int64) bool {
	t := time.Unix(timestamp, 0)

	return (t.Equal(end) || t.After(end)) && (t.Equal(start) || t.Before(start))
}

func chunks[T any](start time.Time, chunkSize int, values []T, at func(T) int64) [][]T {
	if len(values) == 0 {
		return nil
	}

	var result [][]T

	size := int64(chunkSize) * 60
	begin := 0
	current := (start.Unix() - at(values[0])) / size

	for i := 1; i < len(values); i++ {
		bucket := (start.Unix() - at(values[i])) / size

		if bucket != current {
			result = append(result, values[begin:i:i])
			begin = i
			current = bucket
		}
	}

	return append(result, values[begin:])
}

func day(timestamp int64) time.Time {
	t := time.Unix(timestamp, 0)

	return time.Date(
		t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location(),
	)
}

func escape(value string) string {
	var builder strings.Builder

	builder.Grow(len(value))

	for _, r := range value {
		switch r {
		case '*', '?', '[', ']', '\\':
			builder.WriteRune('\\')
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

func findPaths(conn redis.Conn, probe Probe, prefix string, days []int64, scanCount int) ([]string, error) {
	for _, day := range days {
		keyPrefix := string(probe) + ":" + prefix + strconv.FormatInt(day, 10) + ":"

		var paths []string

		cursor := 0

		for {
			values, err := redis.Values(
				conn.Do("SCAN", cursor, "MATCH", keyPrefix+"*", "COUNT", scanCount),
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
			path, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, keyPrefix))
			if err != nil {
				return nil, err
			}

			paths[key] = string(path)
		}

		sort.Strings(paths)

		return paths, nil
	}

	return []string{}, nil
}

func now() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location(),
	)
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
			return [seriesLoadSegmentCount]float64{}, false
		}

		loads[i] = load
	}

	return loads, true
}

func today() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}
