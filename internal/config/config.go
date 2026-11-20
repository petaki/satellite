package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/petaki/satellite/internal/models"
	"github.com/petaki/support-go/cli"
)

// Config type.
type Config struct {
	Debug                  bool
	Name                   string
	Addr                   string
	URL                    string
	AllowedHosts           []string
	SeriesButtons          []models.SeriesType
	RedisURL               string
	MCPEnabled             bool
	HeartbeatEnabled       bool
	HeartbeatWait          int
	HeartbeatSleep         int
	HeartbeatWebhookMethod string
	HeartbeatWebhookURL    string
	HeartbeatWebhookHeader map[string]string
	HeartbeatWebhookBody   string
}

// NewConfig function.
func NewConfig(command *cli.Command, arguments []string) (*Config, error) {
	debug := command.Bool("debug", "APP_DEBUG", false, "Application Debug Mode")
	name := command.String("name", "APP_NAME", "", "Application Name")
	addr := command.String("addr", "APP_ADDR", "", "Application Address")
	appURL := command.String("url", "APP_URL", "", "Application URL")
	allowedHosts := command.String("allowed-hosts", "APP_ALLOWED_HOSTS", "", "Application Allowed Hosts")
	seriesButtons := command.String("series-buttons", "APP_SERIES_BUTTONS", "", "Application Series Buttons")
	redisURL := command.Secret("redis-url", "REDIS_URL", "Redis URL")
	mcpEnabled := command.Bool("mcp-enabled", "MCP_ENABLED", false, "MCP Enabled")
	heartbeatEnabled := command.Bool("heartbeat-enabled", "HEARTBEAT_ENABLED", false, "Heartbeat Enabled")
	heartbeatWait := command.Int("heartbeat-wait", "HEARTBEAT_WAIT", 5, "Heartbeat Wait")
	heartbeatSleep := command.Int("heartbeat-sleep", "HEARTBEAT_SLEEP", 300, "Heartbeat Sleep")
	heartbeatWebhookMethod := command.String("heartbeat-webhook-method", "HEARTBEAT_WEBHOOK_METHOD", "", "Heartbeat Webhook Method")
	heartbeatWebhookURL := command.Secret("heartbeat-webhook-url", "HEARTBEAT_WEBHOOK_URL", "Heartbeat Webhook URL")
	heartbeatWebhookHeader := command.Secret("heartbeat-webhook-header", "HEARTBEAT_WEBHOOK_HEADER", "Heartbeat Webhook Header")
	heartbeatWebhookBody := command.String("heartbeat-webhook-body", "HEARTBEAT_WEBHOOK_BODY", "", "Heartbeat Webhook Body")

	_, err := command.Parse(arguments)
	if err != nil {
		return nil, err
	}

	if *heartbeatWait < 1 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidHeartbeatWait, *heartbeatWait)
	}

	if *heartbeatSleep < 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidHeartbeatSleep, *heartbeatSleep)
	}

	if *heartbeatEnabled && *heartbeatWebhookURL == "" {
		return nil, ErrMissingHeartbeatWebhookURL
	}

	var webhookHeader map[string]string

	if *heartbeatWebhookHeader != "" {
		err = json.Unmarshal([]byte(*heartbeatWebhookHeader), &webhookHeader)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidHeartbeatWebhookHeader, err)
		}
	}

	appConfig := &Config{
		Debug:                  *debug,
		Name:                   *name,
		Addr:                   *addr,
		URL:                    *appURL,
		RedisURL:               *redisURL,
		MCPEnabled:             *mcpEnabled,
		HeartbeatEnabled:       *heartbeatEnabled,
		HeartbeatWait:          *heartbeatWait,
		HeartbeatSleep:         *heartbeatSleep,
		HeartbeatWebhookMethod: *heartbeatWebhookMethod,
		HeartbeatWebhookURL:    *heartbeatWebhookURL,
		HeartbeatWebhookHeader: webhookHeader,
		HeartbeatWebhookBody:   *heartbeatWebhookBody,
	}

	err = appConfig.parseSeriesButtons(*seriesButtons)
	if err != nil {
		return nil, err
	}

	err = appConfig.parseAllowedHosts(*allowedHosts)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

func (c *Config) parseSeriesButtons(value string) error {
	var sb []models.SeriesType

	segments := strings.SplitSeq(value, ",")

	for segment := range segments {
		st := models.SeriesType(strings.TrimSpace(segment))

		for _, current := range models.SeriesTypes {
			if st == current["value"].(models.SeriesType) {
				sb = append(sb, st)

				break
			}
		}
	}

	if len(sb) == 0 {
		return fmt.Errorf("%w: series buttons, %q", ErrInvalid, value)
	}

	c.SeriesButtons = sb[:min(4, len(sb))]

	return nil
}

func (c *Config) parseAllowedHosts(value string) error {
	var hosts []string

	for segment := range strings.SplitSeq(value, ",") {
		host := strings.ToLower(strings.TrimSpace(segment))

		if host != "" {
			hosts = append(hosts, host)
		}
	}

	if len(hosts) > 0 {
		c.AllowedHosts = hosts

		return nil
	}

	u, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("%w: %q, %w", ErrInvalidURL, c.URL, err)
	}

	if u.Hostname() == "" {
		return fmt.Errorf("%w: %q", ErrInvalidURL, c.URL)
	}

	c.AllowedHosts = []string{strings.ToLower(u.Hostname())}

	return nil
}
