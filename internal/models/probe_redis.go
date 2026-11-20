package models

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	heartbeatKeySuffix = ":heartbeat"
	probeRecentDays    = 31
	probeScanCount     = 1000
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
	var names []string

	for {
		values, err := redis.Values(
			conn.Do("SCAN", cursor, "MATCH", "*"+seriesCPUKeyPrefix+"*", "COUNT", probeScanCount),
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

		for _, key := range current {
			name, ok := probeName(key)

			if ok && !slices.Contains(names, name) {
				names = append(names, name)
			}
		}

		if cursor == 0 {
			break
		}
	}

	sort.Strings(names)

	probes := make([]Probe, len(names))

	for key, value := range names {
		probes[key] = Probe(value)
	}

	return probes, nil
}

// Exists function.
func (rpr *RedisProbeRepository) Exists(probe Probe) (bool, error) {
	conn := rpr.RedisPool.Get()
	defer conn.Close()

	keyPrefix := string(probe) + ":" + seriesCPUKeyPrefix
	midnight := today()
	args := redis.Args{}

	for i := range probeRecentDays {
		args = args.Add(keyPrefix + strconv.FormatInt(midnight.AddDate(0, 0, -i).Unix(), 10))
	}

	count, err := redis.Int(conn.Do("EXISTS", args...))
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	keys, err := scanKeys(conn, escape(keyPrefix)+"*", probeScanCount)
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		name, ok := probeName(key)

		if ok && name == string(probe) {
			return true, nil
		}
	}

	return false, nil
}

// FindLatestValues function.
func (rpr *RedisProbeRepository) FindLatestValues(probe Probe, limit int) ([]any, *time.Time, error) {
	if limit < 1 {
		return nil, nil, ErrInvalidLimit
	}

	conn := rpr.RedisPool.Get()
	defer conn.Close()

	days := map[string][]string{}

	end := now()
	start := end.Add(-time.Duration(limit-1) * time.Minute)

	for current := start; !current.After(end); current = current.Add(time.Minute) {
		dayKey := strconv.FormatInt(day(current.Unix()).Unix(), 10)

		days[dayKey] = append(days[dayKey], strconv.FormatInt(current.Unix(), 10))
	}

	var values []any

	for day, fields := range days {
		dayValues, err := redis.Values(
			conn.Do("HMGET", redis.Args{}.Add(string(probe)+":"+seriesCPUKeyPrefix+day).AddFlat(fields)...),
		)
		if err != nil {
			return nil, nil, err
		}

		values = append(values, dayValues...)
	}

	return values, &start, nil
}

// HasHeartbeat function.
func (rpr *RedisProbeRepository) HasHeartbeat(probe Probe) (bool, error) {
	conn := rpr.RedisPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", string(probe)+heartbeatKeySuffix))
}

// SetHeartbeat function.
func (rpr *RedisProbeRepository) SetHeartbeat(probe Probe, sleep int) error {
	conn := rpr.RedisPool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"SET", string(probe)+heartbeatKeySuffix, true,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"EXPIRE", string(probe)+heartbeatKeySuffix, sleep,
	)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// Delete function.
func (rpr *RedisProbeRepository) Delete(probe Probe) error {
	conn := rpr.RedisPool.Get()
	defer conn.Close()

	keys := []string{
		string(probe) + ":" + alarmKeyPrefix,
		string(probe) + heartbeatKeySuffix,
	}

	for _, prefix := range []string{
		seriesCPUKeyPrefix,
		seriesMemoryKeyPrefix,
		seriesProcessCPUKeyPrefix,
		seriesProcessMemoryKeyPrefix,
		seriesLoadKeyPrefix,
		seriesDiskKeyPrefix,
		logKeyPrefix,
		alarmKeyPrefix + ":",
	} {
		current, err := scanKeys(conn, escape(string(probe)+":"+prefix)+"*", probeScanCount)
		if err != nil {
			return err
		}

		keys = append(keys, current...)
	}

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = conn.Send("DEL", key)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
