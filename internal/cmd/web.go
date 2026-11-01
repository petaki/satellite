package cmd

import (
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/web"
	"github.com/petaki/support-go/cli"
)

// App is the CLI application instance.
var App *cli.App

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	appConfig, err := config.NewConfig(command, arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	err = web.Serve(App, appConfig)
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}
