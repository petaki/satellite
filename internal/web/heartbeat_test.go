package web

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
)

type stubHeartbeatRepo struct {
	setHeartbeat bool
}

func (s *stubHeartbeatRepo) FindAll() ([]models.Probe, error) { return nil, nil }
func (s *stubHeartbeatRepo) FindLatestValues(models.Probe, int) ([]any, *time.Time, error) {
	start := time.Now()

	return []any{nil, nil}, &start, nil
}
func (s *stubHeartbeatRepo) HasHeartbeat(models.Probe) (bool, error) { return false, nil }
func (s *stubHeartbeatRepo) SetHeartbeat(models.Probe, int) error {
	s.setHeartbeat = true

	return nil
}
func (s *stubHeartbeatRepo) Delete(models.Probe) error { return nil }

func TestHeartbeatTreatsFourHundredAsFailure(t *testing.T) {
	cases := map[int]bool{
		200: true,
		204: true,
		399: true,
		400: false,
		404: false,
		500: false,
	}

	for status, wantNotified := range cases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
		}))

		repo := &stubHeartbeatRepo{}

		var errorLog bytes.Buffer

		webApp := &app{
			appConfig: &config.Config{
				HeartbeatWait:          2,
				HeartbeatSleep:         300,
				HeartbeatWebhookMethod: http.MethodPost,
				HeartbeatWebhookURL:    server.URL,
				HeartbeatWebhookBody:   `{"probe": "%p"}`,
			},
			infoLog:         log.New(io.Discard, "", 0),
			errorLog:        log.New(&errorLog, "", 0),
			client:          server.Client(),
			probeRepository: repo,
		}

		var wg sync.WaitGroup

		wg.Add(1)
		webApp.handleProbe("web-01", &wg)
		wg.Wait()

		server.Close()

		if repo.setHeartbeat != wantNotified {
			t.Errorf("status %d: notification recorded = %v, want %v", status, repo.setHeartbeat, wantNotified)
		}

		logged := errorLog.String()

		if wantNotified {
			if logged != "" {
				t.Errorf("status %d: logged %q, want nothing", status, logged)
			}

			continue
		}

		for _, want := range []string{ErrBadStatusCode.Error(), "web-01", strconv.Itoa(status)} {
			if !strings.Contains(logged, want) {
				t.Errorf("status %d: logged %q, want it to name %q", status, logged, want)
			}
		}
	}
}
