package main

import (
	"log"
	"os"
	"sort"

	"github.com/dominikus1993/dev-news-bot/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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
				EnvVars:  []string{"DISCORD_WEBHOOK_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "discord-webhook-token",
				Value:    "",
				Usage:    "discord-webhook-token",
				EnvVars:  []string{"DISCORD_WEBHOOK_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "mongo-connection-string",
				Value:    "",
				Usage:    "mongo-connection-string",
				EnvVars:  []string{"MONGO_CONNECTION"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "teams-webhook-url",
				Value:   "",
				Usage:   "teams-webhook-url",
				EnvVars: []string{"TEAMS_WEBHOOK_URL"},
			},
		},
		Action: cmd.Parse,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
