package models

import (
	"encoding/base64"
	"sort"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

const (
	logKeyPrefix = "log:"
	logScanCount = 1000
)

// RedisLogRepository type.
type RedisLogRepository struct {
	RedisPool *redis.Pool
}

// FindLogPaths function.
func (rlr *RedisLogRepository) FindLogPaths(probe Probe) ([]string, error) {
	conn := rlr.RedisPool.Get()
	defer conn.Close()

	todayTS := today()

	return findPaths(conn, probe, logKeyPrefix, []int64{
		todayTS.Unix(),
		todayTS.AddDate(0, 0, -1).Unix(),
	}, logScanCount)
}

// FindLog function.
func (rlr *RedisLogRepository) FindLog(probe Probe, path string) ([]LogEntry, error) {
	conn := rlr.RedisPool.Get()
	defer conn.Close()

	entries := []LogEntry{}

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
