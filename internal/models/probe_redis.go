package models

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/exp/slices"
)

const (
	heartbeatKeySuffix = ":heartbeat"
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
			conn.Do("SCAN", cursor, "MATCH", "*"+seriesCPUKeyPrefix+"*"),
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

		for _, name := range current {
			if strings.Contains(name, seriesProcessCPUKeyPrefix) {
				continue
			}

			segments := strings.SplitN(name, ":"+seriesCPUKeyPrefix, 2)

			if !slices.Contains(names, segments[0]) {
				names = append(names, segments[0])
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

// FindLatestValues function.
func (rpr *RedisProbeRepository) FindLatestValues(probe Probe, limit int) ([]interface{}, *time.Time, error) {
	if limit < 1 {
		return nil, nil, ErrInvalidLimit
	}

	conn := rpr.RedisPool.Get()
	defer conn.Close()

	days := map[string][]string{}

	now := time.Now()
	end := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		now.Location(),
	)

	start := end.Add(-time.Duration(limit-1) * time.Minute)

	for current := start; !current.After(end); current = current.Add(time.Minute) {
		day := strconv.FormatInt(time.Date(
			current.Year(),
			current.Month(),
			current.Day(),
			0,
			0,
			0,
			0,
			now.Location(),
		).Unix(), 10)

		days[day] = append(days[day], strconv.FormatInt(current.Unix(), 10))
	}

	var values []interface{}

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

	cursor := 0
	var keys []string

	for {
		values, err := redis.Values(
			conn.Do("SCAN", cursor, "MATCH", string(probe)+":*"),
		)
		if err != nil {
			return err
		}

		cursor, err = redis.Int(values[0], nil)
		if err != nil {
			return err
		}

		current, err := redis.Strings(values[1], nil)
		if err != nil {
			return err
		}

		keys = append(keys, current...)

		if cursor == 0 {
			break
		}
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
