package cmd

import (
	"context"
	"flag"

	"github.com/dominikus1993/dev-news-bot/src/core/notifications"
	"github.com/dominikus1993/dev-news-bot/src/core/parsers"
	"github.com/dominikus1993/dev-news-bot/src/core/providers"
	"github.com/dominikus1993/dev-news-bot/src/core/usecase"
	"github.com/dominikus1993/dev-news-bot/src/infrastructure/notifiers"
	iparsers "github.com/dominikus1993/dev-news-bot/src/infrastructure/parser"
	irepositories "github.com/dominikus1993/dev-news-bot/src/infrastructure/repositories"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type ParseArticlesAndSendIt struct {
	quantity              int
	dicordWebhookId       string
	discordWebhookToken   string
	mongoConnectionString string
}

func (*ParseArticlesAndSendIt) Name() string     { return "parse" }
func (*ParseArticlesAndSendIt) Synopsis() string { return "Parse Articles And Send It" }
func (*ParseArticlesAndSendIt) Usage() string {
	return `go run . parse"`
}

func (p *ParseArticlesAndSendIt) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.quantity, "quantity", 10, "quantity of articles")
	f.StringVar(&p.dicordWebhookId, "dicord-webhook-id", "", "dicord webhook id")
	f.StringVar(&p.discordWebhookToken, "discord-webhook-token", "", "discord webhook token")
	f.StringVar(&p.mongoConnectionString, "mongo-connection-string", "", "mongo connection string")
}

func (p *ParseArticlesAndSendIt) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log.WithField("mongo", p.mongoConnectionString).WithField("quantity", p.quantity).WithField("dicord-webhook-id", p.dicordWebhookId).WithField("discord-webhook-token", p.discordWebhookToken).Infoln("Parse Articles And Send It")
	mongodbClient, err := irepositories.NewClient(ctx, p.mongoConnectionString)
	if err != nil {
		log.WithError(err).Error("can't create mongodb client")
	}
	defer mongodbClient.Close(ctx)
	redditParser := iparsers.NewRedditParser([]string{"dotnet", "csharp", "fsharp", "golang", "python", "node", "javascript", "devops"})
	hackernewsParser := iparsers.NewHackerNewsArticleParser()
	dotnetomaniakParser := iparsers.NewDotnetoManiakParser()
	parsers := []parsers.ArticlesParser{redditParser, hackernewsParser, dotnetomaniakParser}
	repo := irepositories.NewMongoArticlesRepository(mongodbClient, "Articles")
	articlesProvider := providers.NewArticlesProvider(parsers)
	discord, err := notifiers.NewDiscordWebhookNotifier(p.dicordWebhookId, p.discordWebhookToken)
	if err != nil {
		log.WithError(err).Error("error creating discord notifier")
		return subcommands.ExitFailure
	}
	defer discord.Close()
	bradcaster := notifications.NewBroadcaster([]notifications.Notifier{discord})
	usecase := usecase.NewParseArticlesAndSendItUseCase(articlesProvider, repo, bradcaster)

	err = usecase.Execute(ctx, p.quantity)

	if err != nil {
		log.WithError(err).Error("Error while parsing articles")
		return subcommands.ExitFailure
	}
	log.Infoln("Parsing articles finished")
	return subcommands.ExitSuccess
}
