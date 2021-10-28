package cmd

import (
	"context"
	"flag"

	"github.com/dominikus1993/dev-news-bot/src/core/usecase"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type ParseArticlesAndSendIt struct {
	quantity int
}

func (*ParseArticlesAndSendIt) Name() string     { return "parse" }
func (*ParseArticlesAndSendIt) Synopsis() string { return "Parse Articles And Send It" }
func (*ParseArticlesAndSendIt) Usage() string {
	return `go run . parse"`
}

func (p *ParseArticlesAndSendIt) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.quantity, "quantity", 10, "quantity of articles")
}

func (p *ParseArticlesAndSendIt) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	usecase := usecase.NewParseArticlesAndSendItUseCase()

	err := usecase.Execute(ctx, p.quantity)

	if err != nil {
		log.WithError(err).Error("Error while parsing articles")
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
