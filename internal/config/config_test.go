package config

import (
	"testing"

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
