package cmd

import (
	"context"
	"flag"

	"github.com/dominikus1993/dev-news-bot/internal/console"
	"github.com/dominikus1993/dev-news-bot/internal/discord"
	"github.com/dominikus1993/dev-news-bot/internal/mongo"
	"github.com/dominikus1993/dev-news-bot/internal/parser/devto"
	"github.com/dominikus1993/dev-news-bot/internal/parser/dotnetomaniak"
	"github.com/dominikus1993/dev-news-bot/internal/parser/hackernews"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/usecase"
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
	mongodbClient, err := mongo.NewClient(ctx, p.mongoConnectionString, "Articles")
	if err != nil {
		log.WithError(err).Error("can't create mongodb client")
		return subcommands.ExitFailure
	}
	defer mongodbClient.Close(ctx)
	devtoParser := devto.NewDevToParser([]string{"dotnet", "csharp", "fsharp", "golang", "python", "node", "javascript", "devops", "rust", "aws"})
	hackernewsParser := hackernews.NewHackerNewsArticleParser()
	dotnetomaniakParser := dotnetomaniak.NewDotnetoManiakParser()
	parsers := []parsers.ArticlesParser{hackernewsParser, dotnetomaniakParser, devtoParser}
	repo := mongo.NewMongoArticlesRepository(mongodbClient)
	articlesProvider := providers.NewArticlesProvider(parsers, repo)
	discord, err := discord.NewDiscordWebhookNotifier(p.dicordWebhookId, p.discordWebhookToken)
	if err != nil {
		log.WithError(err).Error("error creating discord notifier")
		return subcommands.ExitFailure
	}
	defer discord.Close()
	pprint := console.NewPPrintNotifier()
	bradcaster := notifications.NewBroadcaster(pprint, discord)
	usecase := usecase.NewParseArticlesAndSendItUseCase(articlesProvider, repo, bradcaster)

	err = usecase.Execute(ctx, p.quantity)

	if err != nil {
		log.WithError(err).Error("Error while parsing articles")
		return subcommands.ExitFailure
	}
	log.Infoln("Parsing articles finished")
	return subcommands.ExitSuccess
}
