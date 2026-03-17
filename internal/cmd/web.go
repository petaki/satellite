package cmd

import (
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/web"
	"github.com/petaki/support-go/cli"
)

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	appConfig, err := config.NewConfig(command, arguments)
	if err != nil {
		cli.ErrorLog.Fatal(err)

		return command.PrintHelp(group)
	}

	web.Serve(appConfig)

	return cli.Success
}
