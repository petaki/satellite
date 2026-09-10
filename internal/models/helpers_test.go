package models

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"

	"github.com/petaki/satellite/internal/service/fake"
)

func TestBetween(t *testing.T) {
	start := time.Date(2026, 9, 10, 12, 0, 0, 0, time.Local)
	end := start.Add(-5 * time.Minute)

	cases := map[string]struct {
		at   time.Time
		want bool
	}{
		"at the newer bound": {start, true},
		"at the older bound": {end, true},
		"inside the window":  {start.Add(-1 * time.Minute), true},
		"newer than start":   {start.Add(time.Second), false},
		"older than end":     {end.Add(-time.Second), false},
	}

	for name, c := range cases {
		if got := between(start, end, c.at.Unix()); got != c.want {
			t.Errorf("%s: between(%s, %s, %s) = %v, want %v",
				name, start.Format(time.TimeOnly), end.Format(time.TimeOnly),
				c.at.Format(time.TimeOnly), got, c.want)
		}
	}
}

func TestChunks(t *testing.T) {
	start := time.Date(2026, 9, 10, 12, 0, 0, 0, time.Local)
	at := func(v int64) int64 { return v }

	minutesAgo := func(m ...int) []int64 {
		values := make([]int64, len(m))

		for i, v := range m {
			values[i] = start.Add(time.Duration(-v) * time.Minute).Unix()
		}

		return values
	}

	sizes := func(result [][]int64) []int {
		out := make([]int, len(result))

		for i, c := range result {
			out[i] = len(c)
		}

		return out
	}

	cases := []struct {
		name      string
		chunkSize int
		values    []int64
		want      []int
	}{
		{"empty", 1, nil, []int{}},
		{"single value", 1, minutesAgo(1), []int{1}},
		{"one minute buckets", 1, minutesAgo(1, 2, 3, 4), []int{1, 1, 1, 1}},
		{"one hour bucket holds all", 60, minutesAgo(1, 2, 3, 4), []int{4}},
		{"one hour bucket splits", 60, minutesAgo(1, 2, 61, 62), []int{2, 2}},
		{"gap wider than the bucket", 60, minutesAgo(1, 200), []int{1, 1}},
		{"bucket boundary", 60, minutesAgo(0, 59, 60), []int{2, 1}},
		{"a day per bucket", 1440, minutesAgo(1, 1441, 2881), []int{1, 1, 1}},
	}

	for _, c := range cases {
		got := sizes(chunks(start, c.chunkSize, c.values, at))

		if len(got) != len(c.want) {
			t.Errorf("%s: got %v chunks, want %v", c.name, got, c.want)

			continue
		}

		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: got %v, want %v", c.name, got, c.want)

				break
			}
		}
	}
}

func TestChunksAreOrderedNewestFirst(t *testing.T) {
	start := time.Date(2026, 9, 10, 12, 0, 0, 0, time.Local)
	at := func(v int64) int64 { return v }

	values := make([]int64, 0, 240)

	for i := range 240 {
		values = append(values, start.Add(time.Duration(-i)*time.Minute).Unix())
	}

	previous := int64(1<<63 - 1)

	for _, chunk := range chunks(start, 60, values, at) {
		if chunk[0] > previous {
			t.Fatalf("chunk starting at %d is newer than the previous chunk", chunk[0])
		}

		previous = chunk[0]
	}
}

func TestChunksKeepsEverySample(t *testing.T) {
	start := time.Date(2026, 9, 10, 12, 0, 0, 0, time.Local)
	at := func(v int64) int64 { return v }

	values := make([]int64, 0, 500)

	for i := range 500 {
		values = append(values, start.Add(time.Duration(-i)*time.Minute).Unix())
	}

	for _, chunkSize := range []int{1, 60, 1440, 5760} {
		total := 0

		for _, chunk := range chunks(start, chunkSize, values, at) {
			total += len(chunk)
		}

		if total != len(values) {
			t.Errorf("chunkSize %d: %d samples in, %d out", chunkSize, len(values), total)
		}
	}
}

func TestDay(t *testing.T) {
	cases := []time.Time{
		time.Date(2026, 9, 10, 0, 0, 0, 0, time.Local),
		time.Date(2026, 9, 10, 12, 34, 56, 789, time.Local),
		time.Date(2026, 9, 10, 23, 59, 59, 0, time.Local),
		time.Date(2026, 3, 1, 2, 30, 0, 0, time.Local),
	}

	for _, at := range cases {
		got := day(at.Unix())

		if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
			t.Errorf("day(%s) = %s, want midnight", at, got)
		}

		if got.Year() != at.Year() || got.Month() != at.Month() || got.Day() != at.Day() {
			t.Errorf("day(%s) = %s, want the same calendar day", at, got)
		}
	}
}

func TestEscape(t *testing.T) {
	cases := map[string]string{
		"*":      `\*`,
		"*:*":    `\*:\*`,
		"web-01": "web-01",
		"a?b":    `a\?b`,
		"a[b]c":  `a\[b\]c`,
		`a\b`:    `a\\b`,
	}

	for in, want := range cases {
		if got := escape(in); got != want {
			t.Errorf("escape(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFindPaths(t *testing.T) {
	server := fake.NewRedis()
	server.HSet("web-01:disk:200:"+encodePath("/data"), "1", "1")
	server.HSet("web-01:disk:200:"+encodePath("/"), "1", "1")
	server.HSet("web-01:disk:100:"+encodePath("/old"), "1", "1")

	paths, err := findPathsWith(server, "web-01", "disk:", []int64{200, 100}, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 2 || paths[0] != "/" || paths[1] != "/data" {
		t.Errorf("paths = %v, want [/ /data] sorted from the first day with data", paths)
	}
}

func TestFindPathsFallsBackToTheNextDay(t *testing.T) {
	server := fake.NewRedis()
	server.HSet("web-01:disk:100:"+encodePath("/old"), "1", "1")

	paths, err := findPathsWith(server, "web-01", "disk:", []int64{200, 100}, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 1 || paths[0] != "/old" {
		t.Errorf("paths = %v, want [/old] from the second day", paths)
	}
}

func TestFindPathsPagesThroughTheCursor(t *testing.T) {
	server := fake.NewRedis()

	for i := range 25 {
		server.HSet("web-01:log:200:"+encodePath("/var/log/"+strconv.Itoa(i)+".log"), "1", "1")
	}

	paths, err := findPathsWith(server, "web-01", "log:", []int64{200}, 4)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 25 {
		t.Errorf("got %d paths, want 25", len(paths))
	}

	if server.Calls("SCAN") < 2 {
		t.Errorf("SCAN was called %d times; the cursor loop was not exercised", server.Calls("SCAN"))
	}
}

func TestFindPathsRejectsUndecodableKeys(t *testing.T) {
	server := fake.NewRedis()
	server.HSet("web-01:disk:200:not!base64", "1", "1")

	if _, err := findPathsWith(server, "web-01", "disk:", []int64{200}, 10); err == nil {
		t.Error("want an error for a key whose suffix is not base64")
	}
}

func TestFindPathsWithoutKeysIsEmptyNotNil(t *testing.T) {
	paths, err := findPathsWith(fake.NewRedis(), "web-01", "disk:", []int64{200, 100}, 10)
	if err != nil {
		t.Fatal(err)
	}

	if paths == nil {
		t.Fatal("paths is nil; it marshals to null instead of []")
	}

	if len(paths) != 0 {
		t.Errorf("paths = %v, want empty", paths)
	}
}

func TestNow(t *testing.T) {
	got := now()

	if got.Second() != 0 || got.Nanosecond() != 0 {
		t.Errorf("now() = %s, want it truncated to the minute", got)
	}

	if delta := time.Since(got); delta < 0 || delta > time.Minute {
		t.Errorf("now() = %s, which is %s away from the wall clock", got, delta)
	}
}

func TestParseLoads(t *testing.T) {
	valid := map[string][seriesLoadSegmentCount]float64{
		"1.500000:2.500000:3.500000": {1.5, 2.5, 3.5},
		"1:2:3":                      {1, 2, 3},
		"0:0:0":                      {0, 0, 0},
		"-1.5:0:2":                   {-1.5, 0, 2},
	}

	for value, want := range valid {
		loads, ok := parseLoads(value)

		if !ok {
			t.Errorf("parseLoads(%q) was rejected", value)

			continue
		}

		if loads != want {
			t.Errorf("parseLoads(%q) = %v, want %v", value, loads, want)
		}
	}

	invalid := []string{"1.5:2.5", "1.5", "", "a:b:c", "1.5:b:3.5", "1:2:3:4", "1 : 2 : 3"}

	for _, value := range invalid {
		loads, ok := parseLoads(value)

		if ok {
			t.Errorf("parseLoads(%q) was accepted", value)

			continue
		}

		if loads != ([seriesLoadSegmentCount]float64{}) {
			t.Errorf("parseLoads(%q) rejected the value but returned %v, want zeroes", value, loads)
		}
	}
}

func TestParseLoadsKeepsSegmentOrder(t *testing.T) {
	loads, ok := parseLoads("1:5:15")
	if !ok {
		t.Fatal("parseLoads rejected a well-formed value")
	}

	if loads[0] != 1 || loads[1] != 5 || loads[2] != 15 {
		t.Errorf("loads = %v, want [1 5 15]; load1/load5/load15 are swapped", loads)
	}
}

func TestTimestampsCoverTheWholeWindow(t *testing.T) {
	rsr := &RedisSeriesRepository{}
	start := now()

	for _, st := range []SeriesType{
		Last5Minutes, Last15Minutes, Last30Minutes, Last1Hour, Last3Hours,
		Last6Hours, Last12Hours, Last24Hours, Last2Days, Last7Days, Last30Days,
	} {
		oldest := rsr.end(start, st)
		needed := day(oldest.Unix()).Unix()

		covered := false

		for _, ts := range rsr.timestamps(st) {
			if ts == needed {
				covered = true
			}
		}

		if !covered {
			t.Errorf("%s: day bucket %d holds the oldest sample in the window but is never queried", st, needed)
		}
	}
}

func TestToday(t *testing.T) {
	got := today()

	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
		t.Errorf("today() = %s, want midnight", got)
	}

	wall := time.Now()

	if got.Year() != wall.Year() || got.Month() != wall.Month() || got.Day() != wall.Day() {
		t.Errorf("today() = %s, want today's calendar day", got)
	}
}

func encodePath(path string) string {
	return base64.StdEncoding.EncodeToString([]byte(path))
}

func findPathsWith(server *fake.Redis, probe Probe, prefix string, days []int64, scanCount int) ([]string, error) {
	conn := server.Pool().Get()
	defer conn.Close()

	return findPaths(conn, probe, prefix, days, scanCount)
}
