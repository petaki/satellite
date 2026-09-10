package models

import (
	"encoding/base64"
	"slices"
	"strconv"
	"testing"

	"github.com/petaki/satellite/internal/service/fake"
)

func logKey(probe string, at int64, path string) string {
	return probe + ":" + logKeyPrefix + strconv.FormatInt(at, 10) + ":" +
		base64.StdEncoding.EncodeToString([]byte(path))
}

func TestFindLogPathsDecodesAndSortsPaths(t *testing.T) {
	server := fake.NewRedis()

	server.HSet(logKey("web-01", today().Unix(), "/var/log/syslog"), "1", "a")
	server.HSet(logKey("web-01", today().Unix(), "/var/log/auth.log"), "1", "b")

	rlr := &RedisLogRepository{RedisPool: server.Pool()}

	paths, err := rlr.FindLogPaths("web-01")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"/var/log/auth.log", "/var/log/syslog"}

	if !slices.Equal(paths, want) {
		t.Errorf("FindLogPaths() = %v, want %v", paths, want)
	}
}

func TestFindLogPathsFallsBackToYesterday(t *testing.T) {
	server := fake.NewRedis()

	yesterday := today().AddDate(0, 0, -1).Unix()
	server.HSet(logKey("web-01", yesterday, "/var/log/syslog"), "1", "a")

	rlr := &RedisLogRepository{RedisPool: server.Pool()}

	paths, err := rlr.FindLogPaths("web-01")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"/var/log/syslog"}

	if !slices.Equal(paths, want) {
		t.Errorf("FindLogPaths() = %v, want %v — yesterday was not consulted", paths, want)
	}
}

func TestFindLogPathsWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rlr := &RedisLogRepository{RedisPool: fake.NewRedis().Pool()}

	paths, err := rlr.FindLogPaths("web-01")
	if err != nil {
		t.Fatal(err)
	}

	if paths == nil {
		t.Fatal("FindLogPaths() = nil, want an empty slice")
	}

	if len(paths) != 0 {
		t.Errorf("FindLogPaths() = %v, want no paths", paths)
	}
}

func TestFindLogReturnsNewestFirst(t *testing.T) {
	server := fake.NewRedis()

	key := logKey("web-01", today().Unix(), "/var/log/syslog")
	server.HSet(key, "100", "oldest")
	server.HSet(key, "300", "newest")
	server.HSet(key, "200", "middle")

	rlr := &RedisLogRepository{RedisPool: server.Pool()}

	entries, err := rlr.FindLog("web-01", "/var/log/syslog")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"newest", "middle", "oldest"}

	if len(entries) != len(want) {
		t.Fatalf("FindLog() = %d entries, want %d", len(entries), len(want))
	}

	for i, content := range want {
		if entries[i].Content != content {
			t.Errorf("entry %d = %q, want %q", i, entries[i].Content, content)
		}
	}
}

func TestFindLogReachesBackWhenTodayStartsAtMidnight(t *testing.T) {
	server := fake.NewRedis()

	todayTS := today().Unix()
	yesterdayTS := today().AddDate(0, 0, -1).Unix()

	server.HSet(logKey("web-01", todayTS, "/var/log/syslog"), strconv.FormatInt(todayTS, 10), "first of today")
	server.HSet(logKey("web-01", yesterdayTS, "/var/log/syslog"), strconv.FormatInt(yesterdayTS+10, 10), "late yesterday")

	rlr := &RedisLogRepository{RedisPool: server.Pool()}

	entries, err := rlr.FindLog("web-01", "/var/log/syslog")
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {
		t.Fatalf("FindLog() = %d entries, want 2 — yesterday was not pulled in", len(entries))
	}

	if entries[0].Content != "first of today" || entries[1].Content != "late yesterday" {
		t.Errorf("FindLog() = %v, want today's entry first", entries)
	}
}

func TestFindLogWithoutDataReturnsEmptyNotNil(t *testing.T) {
	rlr := &RedisLogRepository{RedisPool: fake.NewRedis().Pool()}

	entries, err := rlr.FindLog("web-01", "/var/log/syslog")
	if err != nil {
		t.Fatal(err)
	}

	if entries == nil {
		t.Fatal("FindLog() = nil, want an empty slice — a nil renders as null")
	}

	if len(entries) != 0 {
		t.Errorf("FindLog() = %v, want no entries", entries)
	}
}
