package main

import (
	"context"
	"flag"
	"os"

	"github.com/dominikus1993/dev-news-bot/cmd"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	subcommands.Register(&cmd.ParseArticlesAndSendIt{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
