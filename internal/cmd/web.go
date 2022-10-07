package cmd

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/petaki/satellite/internal/web"
	"github.com/petaki/support-go/cli"
)

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")

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

	heartbeatEnabled := command.FlagSet().Bool("heartbeat-enabled", envHeartbeatEnabled, "Heartbeat Enabled")
	heartbeatWait := command.FlagSet().Int("heartbeat-wait", envHeartbeatWait, "Heartbeat Wait")
	heartbeatSleep := command.FlagSet().Int("heartbeat-sleep", envHeartbeatSleep, "Heartbeat Sleep")
	heartbeatWebhookMethod := command.FlagSet().String("heartbeat-webhook-method", os.Getenv("HEARTBEAT_WEBHOOK_METHOD"), "Heartbeat Webhook Method")
	heartbeatWebhookURL := command.FlagSet().String("heartbeat-webhook-url", os.Getenv("HEARTBEAT_WEBHOOK_URL"), "Heartbeat Webhook URL")
	heartbeatWebhookHeader := command.FlagSet().String("heartbeat-webhook-header", os.Getenv("HEARTBEAT_WEBHOOK_HEADER"), "Heartbeat Webhook Header")
	heartbeatWebhookBody := command.FlagSet().String("heartbeat-webhook-body", os.Getenv("HEARTBEAT_WEBHOOK_BODY"), "Heartbeat Webhook Body")

	_, err = command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	redisPool := newRedisPool(*redisURL)
	defer redisPool.Close()

	var envHeartbeatWebhookHeader map[string]string

	err = json.Unmarshal([]byte(*heartbeatWebhookHeader), &envHeartbeatWebhookHeader)
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	web.Serve(
		*debug,
		*addr,
		*url,
		redisPool,
		*heartbeatEnabled,
		*heartbeatWait,
		*heartbeatSleep,
		*heartbeatWebhookMethod,
		*heartbeatWebhookURL,
		envHeartbeatWebhookHeader,
		*heartbeatWebhookBody,
	)

	return cli.Success
}
