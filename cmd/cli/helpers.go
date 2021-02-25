package main

import (
	"os"

	"github.com/petaki/support-go/cli"
)

func createRedisFlags(command *cli.Command) (*string, *string) {
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")
	redisKeyPrefix := command.FlagSet().String("redis-key-prefix", os.Getenv("REDIS_KEY_PREFIX"), "Redis Key Prefix")

	return redisURL, redisKeyPrefix
}
