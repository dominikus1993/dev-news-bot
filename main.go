package main

import (
	"context"
	"log"
	"os"
	"sort"

	"github.com/dominikus1993/dev-news-bot/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "quantity",
				Value: 10,
				Usage: "quantity",
			},
			&cli.StringFlag{
				Name:     "dicord-webhook-id",
				Value:    "",
				Usage:    "dicord-webhook-id",
				Sources:  cli.EnvVars("DISCORD_WEBHOOK_ID"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "discord-webhook-token",
				Value:    "",
				Usage:    "discord-webhook-token",
				Sources:  cli.EnvVars("DISCORD_WEBHOOK_TOKEN"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "mongo-connection-string",
				Value:    "",
				Usage:    "mongo-connection-string",
				Sources:  cli.EnvVars("MONGO_CONNECTION"),
				Required: true,
			},
			&cli.StringFlag{
				Name:    "teams-webhook-url",
				Value:   "",
				Usage:   "teams-webhook-url",
				Sources: cli.EnvVars("TEAMS_WEBHOOK_URL"),
			},
		},
		Action: cmd.Parse,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
