package models

import (
	"encoding/base64"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/petaki/satellite/internal/service/fake"
)

func seed(t *testing.T, minutes int, withProcesses bool) *fake.Redis {
	t.Helper()

	server := fake.NewRedis()

	for i := 1; i <= minutes; i++ {
		at := now().Add(time.Duration(-i) * time.Minute)
		dayKey := strconv.FormatInt(day(at.Unix()).Unix(), 10)
		field := strconv.FormatInt(at.Unix(), 10)

		server.HSet("web-01:cpu:"+dayKey, field, strconv.Itoa(10+i))
		server.HSet("web-01:load:"+dayKey, field, strconv.Itoa(i)+".0:"+strconv.Itoa(i*2)+".0:"+strconv.Itoa(i*3)+".0")

		if withProcesses {
			server.HSet("web-01:process:cpu:"+dayKey, field,
				"firefox:"+strconv.Itoa(50+i)+".0|node:"+strconv.Itoa(5+i)+".0|zsh:"+strconv.Itoa(i)+".0")
		}
	}

	return server
}

func TestFindCPUAlignsProcessSeriesWithTheAverage(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: seed(t, 4, true).Pool()}

	_, _, avg, p1, p2, p3, err := rsr.FindCPU("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if len(avg) != 4 {
		t.Fatalf("avg = %d points, want 4", len(avg))
	}

	for name, series := range map[string]ProcessSeries{"process1": p1, "process2": p2, "process3": p3} {
		if len(series) != len(avg) {
			t.Fatalf("%s has %d points, avg has %d", name, len(series), len(avg))
		}

		seen := map[int64]bool{}

		for i, value := range series {
			if value.X != avg[i].X {
				t.Errorf("%s[%d].X = %d, avg[%d].X = %d (stack misaligned)", name, i, value.X, i, avg[i].X)
			}

			if seen[value.X] {
				t.Errorf("%s has a duplicate x %d", name, value.X)
			}

			seen[value.X] = true
		}
	}

	if p1[0].Name != "firefox" || p2[0].Name != "node" || p3[0].Name != "zsh" {
		t.Errorf("process names = %q/%q/%q, want firefox/node/zsh", p1[0].Name, p2[0].Name, p3[0].Name)
	}
}

func TestFindCPUFetchesProcessesOncePerDayBucket(t *testing.T) {
	server := seed(t, 4, true)
	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	_, _, avg, _, _, _, err := rsr.FindCPU("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if got := server.Calls("HMGET"); got >= len(avg) {
		t.Errorf("HMGET called %d times for %d points; the per-point fetch is back", got, len(avg))
	}
}

func TestFindCPUWithoutProcessDataStillAligns(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: seed(t, 4, false).Pool()}

	_, _, avg, p1, _, _, err := rsr.FindCPU("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	for i := range avg {
		if p1[i].X != avg[i].X {
			t.Errorf("placeholder[%d].X = %d, want %d", i, p1[i].X, avg[i].X)
		}

		if p1[i].Name != "Not Set" {
			t.Errorf("placeholder[%d].Name = %q, want %q", i, p1[i].Name, "Not Set")
		}
	}
}

func TestFindLoadReadsTheKeyOncePerDay(t *testing.T) {
	server := seed(t, 4, false)
	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	load1, load5, load15, err := rsr.FindLoad("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if len(load1) != 4 || len(load5) != 4 || len(load15) != 4 {
		t.Fatalf("load points = %d/%d/%d, want 4 each", len(load1), len(load5), len(load15))
	}

	if load1[0].Y != 1 || load5[0].Y != 2 || load15[0].Y != 3 {
		t.Errorf("newest load = %v/%v/%v, want 1/2/3 (components swapped?)", load1[0].Y, load5[0].Y, load15[0].Y)
	}

	days := len(rsr.timestamps(Last5Minutes))

	if got := server.Calls("HGETALL"); got > days {
		t.Errorf("HGETALL called %d times for %d day buckets; the three components are not sharing a fetch", got, days)
	}
}

func TestFindCPUWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	mn, mx, avg, p1, p2, p3, err := rsr.FindCPU("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if mn == nil || mx == nil || avg == nil || p1 == nil || p2 == nil || p3 == nil {
		t.Error("a nil series marshals to null instead of []")
	}
}

func TestFindLoadWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	load1, load5, load15, err := rsr.FindLoad("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if load1 == nil || load5 == nil || load15 == nil {
		t.Error("a nil load series marshals to null instead of []")
	}
}

func TestFindCPUKeepsOnlyTheTopProcesses(t *testing.T) {
	server := fake.NewRedis()

	at := now().Add(-time.Minute)
	dayKey := strconv.FormatInt(day(at.Unix()).Unix(), 10)
	field := strconv.FormatInt(at.Unix(), 10)

	server.HSet("web-01:cpu:"+dayKey, field, "10")
	server.HSet("web-01:process:cpu:"+dayKey, field, "a:9.0|b:8.0|c:7.0|d:6.0|e:5.0")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	_, _, _, p1, p2, p3, err := rsr.FindCPU("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if p1[0].Name != "a" || p2[0].Name != "b" || p3[0].Name != "c" {
		t.Errorf("kept %q/%q/%q, want the first three a/b/c", p1[0].Name, p2[0].Name, p3[0].Name)
	}
}

func TestFindDiskPathsDecodesAndSortsPaths(t *testing.T) {
	server := fake.NewRedis()

	dayKey := strconv.FormatInt(today().Unix(), 10)

	for _, path := range []string{"/var", "/"} {
		server.HSet(
			"web-01:"+seriesDiskKeyPrefix+dayKey+":"+base64.StdEncoding.EncodeToString([]byte(path)),
			strconv.FormatInt(today().Unix(), 10), "50",
		)
	}

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	paths, err := rsr.FindDiskPaths("web-01")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"/", "/var"}

	if !slices.Equal(paths, want) {
		t.Errorf("FindDiskPaths() = %v, want %v", paths, want)
	}
}

func TestFindDiskPathsWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	paths, err := rsr.FindDiskPaths("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if paths == nil {
		t.Fatal("FindDiskPaths() = nil, want an empty slice")
	}

	if len(paths) != 0 {
		t.Errorf("FindDiskPaths() = %v, want no paths", paths)
	}
}

func TestFindLatestCPUTakesTheNewestField(t *testing.T) {
	server := fake.NewRedis()

	key := "web-01:" + seriesCPUKeyPrefix + strconv.FormatInt(today().Unix(), 10)
	server.HSet(key, strconv.FormatInt(today().Unix()+60, 10), "11")
	server.HSet(key, strconv.FormatInt(today().Unix()+180, 10), "33")
	server.HSet(key, strconv.FormatInt(today().Unix()+120, 10), "22")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	value, found, err := rsr.FindLatestCPU("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if !found {
		t.Fatal("found = false, want the newest sample")
	}

	if value != 33 {
		t.Errorf("FindLatestCPU() = %v, want 33", value)
	}
}

func TestFindLatestMemoryFallsBackToYesterday(t *testing.T) {
	server := fake.NewRedis()

	yesterday := today().AddDate(0, 0, -1)
	key := "web-01:" + seriesMemoryKeyPrefix + strconv.FormatInt(yesterday.Unix(), 10)
	server.HSet(key, strconv.FormatInt(yesterday.Unix()+60, 10), "64")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	value, found, err := rsr.FindLatestMemory("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if !found || value != 64 {
		t.Errorf("FindLatestMemory() = %v, %v, want 64, true — yesterday was not consulted", value, found)
	}
}

func TestFindLatestMetricWithoutDataIsNotFound(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	_, found, err := rsr.FindLatestCPU("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if found {
		t.Error("found = true for a probe with no samples")
	}
}

func TestFindLatestMetricWithMalformedValueIsNotFound(t *testing.T) {
	server := fake.NewRedis()

	key := "web-01:" + seriesCPUKeyPrefix + strconv.FormatInt(today().Unix(), 10)
	server.HSet(key, strconv.FormatInt(today().Unix()+60, 10), "not a number")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	_, found, err := rsr.FindLatestCPU("web-01")
	if err != nil {
		t.Fatalf("a malformed sample must not surface as an error: %v", err)
	}

	if found {
		t.Error("found = true for an unparseable sample")
	}
}

func TestFindLatestLoadSplitsTheSegments(t *testing.T) {
	server := fake.NewRedis()

	key := "web-01:" + seriesLoadKeyPrefix + strconv.FormatInt(today().Unix(), 10)
	server.HSet(key, strconv.FormatInt(today().Unix()+60, 10), "1.5:2.5:3.5")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	load1, load5, load15, found, err := rsr.FindLatestLoad("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if !found || load1 != 1.5 || load5 != 2.5 || load15 != 3.5 {
		t.Errorf("FindLatestLoad() = %v, %v, %v, %v, want 1.5, 2.5, 3.5, true", load1, load5, load15, found)
	}
}

func TestFindLatestLoadWithMalformedValueIsNotFound(t *testing.T) {
	server := fake.NewRedis()

	key := "web-01:" + seriesLoadKeyPrefix + strconv.FormatInt(today().Unix(), 10)
	server.HSet(key, strconv.FormatInt(today().Unix()+60, 10), "1.5:nope:3.5")

	rsr := &RedisSeriesRepository{RedisPool: server.Pool()}

	load1, _, _, found, err := rsr.FindLatestLoad("web-01")
	if err != nil {
		t.Fatalf("a malformed sample must not surface as an error: %v", err)
	}

	if found || load1 != 0 {
		t.Errorf("FindLatestLoad() = %v, %v, want 0, false — a partial parse leaked", load1, found)
	}
}

func TestChunkSizeGrowsWithTheWindow(t *testing.T) {
	rsr := &RedisSeriesRepository{}

	cases := map[SeriesType]int{
		Last5Minutes: 1,
		Last1Hour:    1,
		Last6Hours:   60,
		Last12Hours:  60,
		Last24Hours:  60,
		Last2Days:    60,
		Last7Days:    60 * 24,
		Last30Days:   60 * 24 * 4,
	}

	for seriesType, want := range cases {
		if got := rsr.ChunkSize(seriesType); got != want {
			t.Errorf("ChunkSize(%s) = %d, want %d", seriesType, got, want)
		}
	}
}

func TestFindMemoryWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	minSeries, maxSeries, avg, p1, p2, p3, err := rsr.FindMemory("web-01", Last5Minutes)
	if err != nil {
		t.Fatal(err)
	}

	if minSeries == nil || maxSeries == nil || avg == nil {
		t.Error("a nil series renders as null, want empty")
	}

	if p1 == nil || p2 == nil || p3 == nil {
		t.Error("a nil process series renders as null, want empty")
	}
}

func TestFindDiskWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rsr := &RedisSeriesRepository{RedisPool: fake.NewRedis().Pool()}

	minSeries, maxSeries, avg, err := rsr.FindDisk("web-01", Last5Minutes, "/")
	if err != nil {
		t.Fatal(err)
	}

	if minSeries == nil || maxSeries == nil || avg == nil {
		t.Error("a nil series renders as null, want empty")
	}
}
