package main

import (
	"os"

	"github.com/dominikus1993/dev-news-bot/cmd"
	log "github.com/sirupsen/logrus"
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

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
