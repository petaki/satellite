package web

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/petaki/satellite/internal/models"
)

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

	values, start, err := a.probeRepository.FindLatestValues(probe, a.heartbeatWait)
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

	if a.heartbeatSleep > 0 {
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

	body := strings.ReplaceAll(a.heartbeatWebhookBody, "%p", string(probe))
	body = strings.ReplaceAll(body, "%t", start.Format(time.RFC3339))
	body = strings.ReplaceAll(body, "%x", strconv.FormatInt(start.Unix(), 10))
	body = strings.ReplaceAll(body, "%l", fmt.Sprintf("/cpu?probe=%s", probe))

	req, err := http.NewRequest(a.heartbeatWebhookMethod, a.heartbeatWebhookURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		a.errorLog.Print(err)

		return
	}

	for key, value := range a.heartbeatWebhookHeader {
		req.Header.Set(key, value)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.errorLog.Print(err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		a.errorLog.Print(errors.New("heartbeat: bad status code"))

		return
	}

	if a.heartbeatSleep > 0 {
		err = a.probeRepository.SetHeartbeat(probe, a.heartbeatSleep)
		if err != nil {
			a.errorLog.Print(err)

			return
		}
	}
}
