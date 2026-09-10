package config

import (
	"errors"
	"slices"
	"testing"

	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/cli"
)

func newTestConfig(t *testing.T, header string) (*Config, error) {
	t.Helper()
	t.Setenv("APP_SERIES_BUTTONS", "last_5_minutes,last_1_hour")
	t.Setenv("HEARTBEAT_WEBHOOK_HEADER", header)

	return NewConfig(&cli.Command{Name: "serve"}, nil)
}

func TestConfigWithoutWebhookHeader(t *testing.T) {
	appConfig, err := newTestConfig(t, "")
	if err != nil {
		t.Fatalf("an unset HEARTBEAT_WEBHOOK_HEADER must not stop startup: %v", err)
	}

	if appConfig.HeartbeatWebhookHeader != nil {
		t.Errorf("header map = %v, want nil", appConfig.HeartbeatWebhookHeader)
	}
}

func TestConfigWithWebhookHeader(t *testing.T) {
	appConfig, err := newTestConfig(t, `{"Authorization": "Bearer TOKEN"}`)
	if err != nil {
		t.Fatal(err)
	}

	if appConfig.HeartbeatWebhookHeader["Authorization"] != "Bearer TOKEN" {
		t.Errorf("header = %v", appConfig.HeartbeatWebhookHeader)
	}
}

func TestConfigRejectsMalformedWebhookHeader(t *testing.T) {
	if _, err := newTestConfig(t, "{not json"); err == nil {
		t.Error("malformed HEARTBEAT_WEBHOOK_HEADER should be an error")
	}
}

func TestConfigRejectsInvalidSeriesButtons(t *testing.T) {
	t.Setenv("APP_SERIES_BUTTONS", "nope")

	_, err := NewConfig(&cli.Command{Name: "serve"}, nil)

	if !errors.Is(err, ErrInvalid) {
		t.Errorf("error = %v, want ErrInvalid", err)
	}
}

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
		var appConfig Config

		err := appConfig.parseSeriesButtons(c.value)

		if c.want == nil {
			if !errors.Is(err, ErrInvalid) {
				t.Errorf("%s: parseSeriesButtons(%q) error = %v, want ErrInvalid", name, c.value, err)
			}

			continue
		}

		if err != nil {
			t.Errorf("%s: parseSeriesButtons(%q) = %v", name, c.value, err)

			continue
		}

		if !slices.Equal(appConfig.SeriesButtons, c.want) {
			t.Errorf("%s: parseSeriesButtons(%q) = %v, want %v", name, c.value, appConfig.SeriesButtons, c.want)
		}
	}
}

func TestParseSeriesButtonsAcceptsEverySeriesType(t *testing.T) {
	for _, current := range models.SeriesTypes {
		value := current["value"].(models.SeriesType)

		var appConfig Config

		err := appConfig.parseSeriesButtons(string(value))
		if err != nil {
			t.Errorf("parseSeriesButtons(%q) = %v", value, err)

			continue
		}

		if len(appConfig.SeriesButtons) != 1 || appConfig.SeriesButtons[0] != value {
			t.Errorf("parseSeriesButtons(%q) = %v, want [%v]", value, appConfig.SeriesButtons, value)
		}
	}
}
