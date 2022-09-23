package cmd

import (
	"os"

	"github.com/petaki/satellite/internal/web"
	"github.com/petaki/support-go/cli"
)

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	redisPool := newRedisPool(*redisURL)
	defer redisPool.Close()

	web.Serve(
		*debug,
		*addr,
		*url,
		redisPool,
	)

	return cli.Success
}
