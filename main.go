package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/satellite/internal/cmd"
	"github.com/petaki/support-go/cli"
)

func main() {
	(&cli.App{
		Name:    "Satellite",
		Version: "master",
		Groups: []*cli.Group{
			{
				Name:  "web",
				Usage: "Web commands",
				Commands: []*cli.Command{
					{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: cmd.WebServe,
					},
				},
			},
		},
	}).Execute()
}
