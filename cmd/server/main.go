package main

import (
	"os"
	"restapi/internal/db/migration"
	"restapi/internal/logger"
	"restapi/internal/server"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Blog API"
	app.Description = "Implementing back-end services for blog application"

	app.Commands = []*cli.Command{
		{
			Name:        "migrations",
			Description: "migrations looks at the currently active migration version and will migrate all the way up (applying all up migrations)",
			Action: func(c *cli.Context) error {
				return migration.Up()
			},
		},
		{
			Name:        "drop",
			Description: "drop deletes everything in the database",
			Action: func(c *cli.Context) error {
				return migration.Drop()
			},
		},
		{
			Name:        "start",
			Description: "start the server",
			Action: func(c *cli.Context) error {
				return server.Start()
			},
		},
		{
			Name:        "launch",
			Description: "launch migrate all the way up (applying all up migrations) and start the server",
			Action: func(c *cli.Context) error {
				err := migration.Up()
				if err != nil {
					return err
				}
				return server.Start()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Log().Fatal().Err(err).Msg("failed to run server")
	}
}
