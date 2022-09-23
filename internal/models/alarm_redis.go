package models

import "github.com/gomodule/redigo/redis"

const (
	alarmKeyPrefix = "alarm"
)

// RedisAlarmRepository type.
type RedisAlarmRepository struct {
	RedisPool *redis.Pool
}

// Find function.
func (rar *RedisAlarmRepository) Find(probe Probe) (*Alarm, error) {
	conn := rar.RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", string(probe)+":"+alarmKeyPrefix))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrNoRecord
	}

	var alarm Alarm

	err = redis.ScanStruct(values, &alarm)
	if err != nil {
		return nil, err
	}

	return &alarm, nil
}
