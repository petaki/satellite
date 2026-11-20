package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/petaki/satellite/internal/models"
)

func (a *app) runHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.heartbeat()
		}
	}
}

func (a *app) heartbeat() {
	err := a.handleProbes()
	if err != nil {
		a.errorLog.Print(err)
	}
}

func (a *app) handleProbes() error {
	probes, err := a.probeRepository.FindAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(len(probes))

	for _, probe := range probes {
		go a.handleProbe(probe, &wg)
	}

	wg.Wait()

	return nil
}

func (a *app) handleProbe(probe models.Probe, wg *sync.WaitGroup) {
	defer wg.Done()

	sendWebhook := true

	values, start, err := a.probeRepository.FindLatestValues(probe, a.appConfig.HeartbeatWait)
	if err != nil {
		a.errorLog.Print(err)

		return
	}

	for _, value := range values {
		if value != nil {
			sendWebhook = false

			break
		}
	}

	if !sendWebhook {
		return
	}

	if a.appConfig.HeartbeatSleep > 0 {
		hasHeartbeat, err := a.probeRepository.HasHeartbeat(probe)
		if err != nil {
			a.errorLog.Print(err)

			return
		}

		if hasHeartbeat {
			return
		}
	}

	a.infoLog.Printf("Calling the heartbeat webhook URL for %s...", probe)

	req, err := http.NewRequest(a.appConfig.HeartbeatWebhookMethod, a.appConfig.HeartbeatWebhookURL, strings.NewReader(a.heartbeatBody(probe, *start)))
	if err != nil {
		a.errorLog.Print(fmt.Errorf("%w: %s, %w", ErrWebhook, probe, withoutURL(err)))

		return
	}

	for key, value := range a.appConfig.HeartbeatWebhookHeader {
		req.Header.Set(key, value)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.errorLog.Print(fmt.Errorf("%w: %s, %w", ErrWebhook, probe, withoutURL(err)))

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		a.errorLog.Print(fmt.Errorf("%w: %s, %d", ErrBadStatusCode, probe, resp.StatusCode))

		return
	}

	if a.appConfig.HeartbeatSleep > 0 {
		err = a.probeRepository.SetHeartbeat(probe, a.appConfig.HeartbeatSleep)
		if err != nil {
			a.errorLog.Print(err)

			return
		}
	}
}

func (a *app) heartbeatBody(probe models.Probe, start time.Time) string {
	replace := func(name string) string {
		return strings.NewReplacer(
			"%p", name,
			"%t", start.Format(time.RFC3339),
			"%x", strconv.FormatInt(start.Unix(), 10),
			"%l", "/cpu?probe="+url.QueryEscape(string(probe)),
		).Replace(a.appConfig.HeartbeatWebhookBody)
	}

	body := replace(jsonEscape(string(probe)))
	if json.Valid([]byte(body)) {
		return body
	}

	return replace(string(probe))
}
