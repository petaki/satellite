package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
)

type offlineProbeRepository struct {
	models.ProbeRepository
	start time.Time
}

func (r *offlineProbeRepository) FindLatestValues(models.Probe, int) ([]any, *time.Time, error) {
	return []any{nil}, &r.start, nil
}

func sendHeartbeat(t *testing.T, probe models.Probe, body string) string {
	t.Helper()

	received := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, _ := io.ReadAll(r.Body)
		received <- string(payload)
	}))
	defer server.Close()

	a := &app{
		appConfig: &config.Config{
			HeartbeatWebhookMethod: http.MethodPost,
			HeartbeatWebhookURL:    server.URL,
			HeartbeatWebhookBody:   body,
		},
		infoLog:         log.New(io.Discard, "", 0),
		errorLog:        log.New(io.Discard, "", 0),
		client:          server.Client(),
		probeRepository: &offlineProbeRepository{start: time.Unix(1700000000, 0).UTC()},
	}

	var wg sync.WaitGroup

	wg.Add(1)
	a.handleProbe(probe, &wg)

	return <-received
}

func TestHeartbeatJSONBodyEscapesProbeName(t *testing.T) {
	template := `{"probe": "%p", "timestamp_rfc3339": "%t", "timestamp_unix": %x, "link": "%l"}`

	for _, probe := range []models.Probe{"web-1", `bad"probe`, `back\slash`, "tab\tname"} {
		body := sendHeartbeat(t, probe, template)

		var payload struct {
			Probe     string `json:"probe"`
			Timestamp string `json:"timestamp_rfc3339"`
			Unix      int64  `json:"timestamp_unix"`
			Link      string `json:"link"`
		}

		err := json.Unmarshal([]byte(body), &payload)
		if err != nil {
			t.Fatalf("probe %q: body %s is not valid JSON: %v", probe, body, err)
		}

		if payload.Probe != string(probe) || payload.Timestamp != "2023-11-14T22:13:20Z" || payload.Unix != 1700000000 {
			t.Errorf("probe %q: got %+v", probe, payload)
		}
	}
}

func TestHeartbeatNonJSONBodyKeepsRawValues(t *testing.T) {
	got := sendHeartbeat(t, `bad"probe`, "probe %p went offline at %x")
	want := `probe bad"probe went offline at 1700000000`

	if got != want {
		t.Errorf("body %q, want %q", got, want)
	}
}
