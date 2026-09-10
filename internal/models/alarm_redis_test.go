package models

import (
	"errors"
	"testing"

	"github.com/petaki/satellite/internal/service/fake"
)

func TestAlarmFindWithoutRecord(t *testing.T) {
	rar := &RedisAlarmRepository{RedisPool: fake.NewRedis().Pool()}

	alarm, err := rar.Find("web-01")

	if !errors.Is(err, ErrNoRecord) {
		t.Errorf("err = %v, want ErrNoRecord", err)
	}

	if alarm != nil {
		t.Errorf("alarm = %v, want nil", alarm)
	}
}

func TestAlarmFindReadsEveryThreshold(t *testing.T) {
	server := fake.NewRedis()

	key := "web-01:" + alarmKeyPrefix
	server.HSet(key, "cpu", "80.5")
	server.HSet(key, "memory", "90")
	server.HSet(key, "disk", "70")
	server.HSet(key, "load", "4.25")

	rar := &RedisAlarmRepository{RedisPool: server.Pool()}

	alarm, err := rar.Find("web-01")
	if err != nil {
		t.Fatal(err)
	}

	want := Alarm{CPU: 80.5, Memory: 90, Disk: 70, Load: 4.25}

	if *alarm != want {
		t.Errorf("Find() = %+v, want %+v", *alarm, want)
	}
}
