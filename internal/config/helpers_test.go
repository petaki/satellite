package config

import (
	"slices"
	"testing"

	"github.com/petaki/satellite/internal/models"
)

func TestParseSeriesButtons(t *testing.T) {
	cases := map[string]struct {
		value string
		want  []models.SeriesType
	}{
		"a single button": {
			"last_5_minutes",
			[]models.SeriesType{models.Last5Minutes},
		},
		"order is kept": {
			"last_7_days,last_5_minutes,last_1_hour",
			[]models.SeriesType{models.Last7Days, models.Last5Minutes, models.Last1Hour},
		},
		"unknown names are dropped": {
			"nope,last_1_hour,also_nope",
			[]models.SeriesType{models.Last1Hour},
		},
		"capped at four": {
			"last_5_minutes,last_15_minutes,last_30_minutes,last_1_hour,last_3_hours",
			[]models.SeriesType{models.Last5Minutes, models.Last15Minutes, models.Last30Minutes, models.Last1Hour},
		},
		"duplicates are kept": {
			"last_1_hour,last_1_hour",
			[]models.SeriesType{models.Last1Hour, models.Last1Hour},
		},
		"padding is trimmed": {
			"last_5_minutes, last_1_hour",
			[]models.SeriesType{models.Last5Minutes, models.Last1Hour},
		},
		"tabs and newlines are trimmed": {
			"\tlast_5_minutes ,\nlast_1_hour\t",
			[]models.SeriesType{models.Last5Minutes, models.Last1Hour},
		},
		"padding alone is still nothing": {
			"  ,  ",
			nil,
		},
		"empty": {
			"",
			nil,
		},
		"nothing valid": {
			"nope,also_nope",
			nil,
		},
	}

	for name, c := range cases {
		got := parseSeriesButtons(c.value)

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: parseSeriesButtons(%q) = %v, want %v", name, c.value, got, c.want)
		}
	}
}

func TestParseSeriesButtonsAcceptsEverySeriesType(t *testing.T) {
	for _, current := range models.SeriesTypes {
		value := current["value"].(models.SeriesType)

		got := parseSeriesButtons(string(value))

		if len(got) != 1 || got[0] != value {
			t.Errorf("parseSeriesButtons(%q) = %v, want [%v]", value, got, value)
		}
	}
}
