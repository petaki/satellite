package models

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/petaki/satellite/internal/service/fake"
)

func TestFindAllIgnoresProcessKeysAndDeduplicatesDays(t *testing.T) {
	server := fake.NewRedis()

	todayKey := strconv.FormatInt(today().Unix(), 10)
	yesterdayKey := strconv.FormatInt(today().AddDate(0, 0, -1).Unix(), 10)

	server.HSet("web-01:"+seriesCPUKeyPrefix+todayKey, "1", "10")
	server.HSet("web-01:"+seriesCPUKeyPrefix+yesterdayKey, "1", "10")
	server.HSet("web-01:"+seriesProcessCPUKeyPrefix+todayKey, "1", "firefox:1.0")
	server.HSet("db-01:"+seriesCPUKeyPrefix+todayKey, "1", "10")
	server.HSet("web-01:"+seriesMemoryKeyPrefix+todayKey, "1", "10")

	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	probes, err := rpr.FindAll()
	if err != nil {
		t.Fatal(err)
	}

	want := []Probe{"db-01", "web-01"}

	if !slices.Equal(probes, want) {
		t.Errorf("FindAll() = %v, want %v", probes, want)
	}
}

func TestFindAllWithoutProbesReturnsEmptyNotNil(t *testing.T) {
	rpr := &RedisProbeRepository{RedisPool: fake.NewRedis().Pool()}

	probes, err := rpr.FindAll()
	if err != nil {
		t.Fatal(err)
	}

	if probes == nil {
		t.Fatal("FindAll() = nil, want an empty slice")
	}

	if len(probes) != 0 {
		t.Errorf("FindAll() = %v, want no probes", probes)
	}
}

func TestFindLatestValuesRejectsNonPositiveLimit(t *testing.T) {
	rpr := &RedisProbeRepository{RedisPool: fake.NewRedis().Pool()}

	for _, limit := range []int{0, -1} {
		_, _, err := rpr.FindLatestValues("web-01", limit)

		if err != ErrInvalidLimit {
			t.Errorf("limit %d: err = %v, want ErrInvalidLimit", limit, err)
		}
	}
}

func TestFindLatestValuesReturnsOneSlotPerMinute(t *testing.T) {
	server := fake.NewRedis()

	for i := 1; i <= 2; i++ {
		at := now().Add(time.Duration(-i) * time.Minute)
		dayKey := strconv.FormatInt(day(at.Unix()).Unix(), 10)

		server.HSet("web-01:"+seriesCPUKeyPrefix+dayKey, strconv.FormatInt(at.Unix(), 10), "10")
	}

	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	values, start, err := rpr.FindLatestValues("web-01", 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(values) != 3 {
		t.Fatalf("values = %d, want 3 (one per minute)", len(values))
	}

	var set int

	for _, value := range values {
		if value != nil {
			set++
		}
	}

	if set != 2 {
		t.Errorf("%d of 3 minutes carried a value, want 2", set)
	}

	if start == nil {
		t.Fatal("start = nil, want the first minute of the window")
	}
}

func TestHasHeartbeatReflectsSetHeartbeat(t *testing.T) {
	server := fake.NewRedis()
	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	has, err := rpr.HasHeartbeat("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if has {
		t.Error("HasHeartbeat() = true before any heartbeat was set")
	}

	err = rpr.SetHeartbeat("web-01", 300)
	if err != nil {
		t.Fatal(err)
	}

	has, err = rpr.HasHeartbeat("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if !has {
		t.Error("HasHeartbeat() = false after SetHeartbeat")
	}
}

func TestSetHeartbeatExpiresWithTheConfiguredSleep(t *testing.T) {
	server := fake.NewRedis()
	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	err := rpr.SetHeartbeat("web-01", 300)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := server.Get("web-01" + heartbeatKeySuffix); !ok {
		t.Fatalf("heartbeat key not set, keys = %v", server.Keys())
	}

	if ttl := server.TTL("web-01" + heartbeatKeySuffix); ttl != 300 {
		t.Errorf("ttl = %d, want 300 — without it the suppression never lifts", ttl)
	}
}

func TestDeleteRemovesOnlyTheNamedProbe(t *testing.T) {
	server := fake.NewRedis()

	dayKey := strconv.FormatInt(today().Unix(), 10)

	for _, probe := range []string{"web-01", "web-011", "web*01", "db-01"} {
		server.HSet(probe+":"+seriesCPUKeyPrefix+dayKey, "1", "10")
	}

	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	err := rpr.SetHeartbeat("web-01", 300)
	if err != nil {
		t.Fatal(err)
	}

	err = rpr.Delete("web-01")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"db-01:" + seriesCPUKeyPrefix + dayKey,
		"web*01:" + seriesCPUKeyPrefix + dayKey,
		"web-011:" + seriesCPUKeyPrefix + dayKey,
	}

	if got := server.Keys(); !slices.Equal(got, want) {
		t.Errorf("keys after Delete(web-01) = %v, want %v", got, want)
	}
}

func TestDeleteTreatsGlobCharactersLiterally(t *testing.T) {
	server := fake.NewRedis()

	dayKey := strconv.FormatInt(today().Unix(), 10)

	for _, probe := range []string{"web-01", "web*01", "db-01"} {
		server.HSet(probe+":"+seriesCPUKeyPrefix+dayKey, "1", "10")
	}

	rpr := &RedisProbeRepository{RedisPool: server.Pool()}

	err := rpr.Delete("web*01")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"db-01:" + seriesCPUKeyPrefix + dayKey,
		"web-01:" + seriesCPUKeyPrefix + dayKey,
	}

	if got := server.Keys(); !slices.Equal(got, want) {
		t.Errorf("keys after Delete(web*01) = %v, want %v — the glob was not escaped", got, want)
	}
}
