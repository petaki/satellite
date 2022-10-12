package web

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	for _, probe := range probes {
		sendWebhook := true

		values, start, err := a.probeRepository.FindLatestValues(probe, a.heartbeatWait)
		if err != nil {
			return err
		}

		for _, value := range values {
			if value != nil {
				sendWebhook = false

				break
			}
		}

		if !sendWebhook {
			return nil
		}

		if a.heartbeatSleep > 0 {
			hasHeartbeat, err := a.probeRepository.HasHeartbeat(probe)
			if err != nil {
				return err
			}

			if hasHeartbeat {
				return nil
			}
		}

		a.infoLog.Printf("Calling the heartbeat webhook URL for %s...", probe)

		body := strings.ReplaceAll(a.heartbeatWebhookBody, "%p", string(probe))
		body = strings.ReplaceAll(body, "%t", start.Format(time.RFC3339))
		body = strings.ReplaceAll(body, "%x", strconv.FormatInt(start.Unix(), 10))
		body = strings.ReplaceAll(body, "%l", fmt.Sprintf("/cpu?probe=%s", probe))

		req, err := http.NewRequest(a.heartbeatWebhookMethod, a.heartbeatWebhookURL, bytes.NewBuffer([]byte(body)))
		if err != nil {
			return err
		}

		for key, value := range a.heartbeatWebhookHeader {
			req.Header.Set(key, value)
		}

		resp, err := a.client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode > 400 {
			return errors.New("heartbeat: bad status code")
		}

		if a.heartbeatSleep > 0 {
			err = a.probeRepository.SetHeartbeat(probe, a.heartbeatSleep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
