package models

import (
	"encoding/base64"
	"sort"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	logKeyPrefix = "log:"
)

// RedisLogRepository type.
type RedisLogRepository struct {
	RedisPool *redis.Pool
}

// FindLogPaths function.
func (rlr *RedisLogRepository) FindLogPaths(probe Probe) ([]string, error) {
	conn := rlr.RedisPool.Get()
	defer conn.Close()

	var paths []string

	todayTS := today()
	yesterdayTS := todayTS.AddDate(0, 0, -1)

	for _, ts := range []int64{todayTS.Unix(), yesterdayTS.Unix()} {
		cursor := 0
		prefix := string(probe) + ":" + logKeyPrefix + strconv.FormatInt(ts, 10) + ":"

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

// FindLog function.
func (rlr *RedisLogRepository) FindLog(probe Probe, path string) ([]LogEntry, error) {
	conn := rlr.RedisPool.Get()
	defer conn.Close()

	var entries []LogEntry

	todayTS := today()
	encodedPath := base64.StdEncoding.EncodeToString([]byte(path))

	key := string(probe) + ":" + logKeyPrefix + strconv.FormatInt(todayTS.Unix(), 10) + ":" + encodedPath

	values, err := redis.Strings(
		conn.Do("HGETALL", key),
	)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(values); i += 2 {
		timestamp, err := strconv.ParseInt(values[i], 10, 64)
		if err != nil {
			return nil, err
		}

		entries = append(entries, LogEntry{
			Timestamp: timestamp,
			Content:   values[i+1],
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Timestamp > entries[j].Timestamp
	})

	if len(entries) > 0 && entries[len(entries)-1].Timestamp == todayTS.Unix() {
		yesterdayTS := todayTS.AddDate(0, 0, -1)
		key = string(probe) + ":" + logKeyPrefix + strconv.FormatInt(yesterdayTS.Unix(), 10) + ":" + encodedPath

		values, err = redis.Strings(
			conn.Do("HGETALL", key),
		)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(values); i += 2 {
			timestamp, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return nil, err
			}

			entries = append(entries, LogEntry{
				Timestamp: timestamp,
				Content:   values[i+1],
			})
		}

		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Timestamp > entries[j].Timestamp
		})
	}

	return entries, nil
}
