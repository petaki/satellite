package models

import (
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
