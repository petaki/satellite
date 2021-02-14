package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/support-go/cli"
)

func main() {
	(&cli.App{
		Name:    "Satellite",
		Version: "1.0.0",
		Groups: []*cli.Group{
			&cli.Group{
				Name:  "web",
				Usage: "Web commands",
				Commands: []*cli.Command{
					&cli.Command{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: webServe,
					},
				},
			},
		},
	}).Execute()
}
