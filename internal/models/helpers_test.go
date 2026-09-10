package models

import "testing"

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
