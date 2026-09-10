package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	SeriesButtons          []models.SeriesType
	RedisURL               string
	MCPEnabled             bool
	MCPAllowAnyHost        bool
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
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	name := command.FlagSet().String("name", os.Getenv("APP_NAME"), "Application Name")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	seriesButtons := command.FlagSet().String("series-buttons", os.Getenv("APP_SERIES_BUTTONS"), "Application Series Buttons")
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	envMCPEnabled, err := strconv.ParseBool(os.Getenv("MCP_ENABLED"))
	if err != nil {
		envMCPEnabled = false
	}

	envMCPAllowAnyHost, err := strconv.ParseBool(os.Getenv("MCP_ALLOW_ANY_HOST"))
	if err != nil {
		envMCPAllowAnyHost = false
	}

	envHeartbeatEnabled, err := strconv.ParseBool(os.Getenv("HEARTBEAT_ENABLED"))
	if err != nil {
		envHeartbeatEnabled = false
	}

	envHeartbeatWait, err := strconv.Atoi(os.Getenv("HEARTBEAT_WAIT"))
	if err != nil || envHeartbeatWait < 1 {
		envHeartbeatWait = 5
	}

	envHeartbeatSleep, err := strconv.Atoi(os.Getenv("HEARTBEAT_SLEEP"))
	if err != nil || envHeartbeatSleep < 0 {
		envHeartbeatSleep = 300
	}

	mcpEnabled := command.FlagSet().Bool("mcp-enabled", envMCPEnabled, "MCP Enabled")
	mcpAllowAnyHost := command.FlagSet().Bool("mcp-allow-any-host", envMCPAllowAnyHost, "MCP Allow Any Host")
	heartbeatEnabled := command.FlagSet().Bool("heartbeat-enabled", envHeartbeatEnabled, "Heartbeat Enabled")
	heartbeatWait := command.FlagSet().Int("heartbeat-wait", envHeartbeatWait, "Heartbeat Wait")
	heartbeatSleep := command.FlagSet().Int("heartbeat-sleep", envHeartbeatSleep, "Heartbeat Sleep")
	heartbeatWebhookMethod := command.FlagSet().String("heartbeat-webhook-method", os.Getenv("HEARTBEAT_WEBHOOK_METHOD"), "Heartbeat Webhook Method")
	heartbeatWebhookURL := command.FlagSet().String("heartbeat-webhook-url", os.Getenv("HEARTBEAT_WEBHOOK_URL"), "Heartbeat Webhook URL")
	heartbeatWebhookHeader := command.FlagSet().String("heartbeat-webhook-header", os.Getenv("HEARTBEAT_WEBHOOK_HEADER"), "Heartbeat Webhook Header")
	heartbeatWebhookBody := command.FlagSet().String("heartbeat-webhook-body", os.Getenv("HEARTBEAT_WEBHOOK_BODY"), "Heartbeat Webhook Body")

	_, err = command.Parse(arguments)
	if err != nil {
		return nil, err
	}

	var webhookHeader map[string]string

	if *heartbeatWebhookHeader != "" {
		err = json.Unmarshal([]byte(*heartbeatWebhookHeader), &webhookHeader)
		if err != nil {
			return nil, err
		}
	}

	appConfig := &Config{
		Debug:                  *debug,
		Name:                   *name,
		Addr:                   *addr,
		URL:                    *url,
		RedisURL:               *redisURL,
		MCPEnabled:             *mcpEnabled,
		MCPAllowAnyHost:        *mcpAllowAnyHost,
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
