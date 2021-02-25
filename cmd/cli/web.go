package main

import (
	"os"

	"github.com/petaki/satellite/internal/web"
	"github.com/petaki/support-go/cli"
)

func webServe(group *cli.Group, command *cli.Command, arguments []string) int {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")

	redisURL, redisKeyPrefix := createRedisFlags(command)

	web.Serve(
		*debug,
		*addr,
		*url,
		*redisURL,
		*redisKeyPrefix,
	)

	return 0
}
