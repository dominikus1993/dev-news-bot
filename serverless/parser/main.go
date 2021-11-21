package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	handler "github.com/dominikus1993/dev-news-bot/src/lambda"

	"github.com/dominikus1993/dev-news-bot/src/core/notifications"
	"github.com/dominikus1993/dev-news-bot/src/core/parsers"
	"github.com/dominikus1993/dev-news-bot/src/core/providers"
	"github.com/dominikus1993/dev-news-bot/src/core/usecase"
	"github.com/dominikus1993/dev-news-bot/src/infrastructure/notifiers"
	iparsers "github.com/dominikus1993/dev-news-bot/src/infrastructure/parser"
	irepositories "github.com/dominikus1993/dev-news-bot/src/infrastructure/repositories"
	log "github.com/sirupsen/logrus"
)

var (
	dicordWebhookId, mongoConnectionString, discordWebhookToken string
	parser                                                      *handler.ParseAndSendArticlesHandler
	discordWebhookNotifier                                      *notifiers.DiscordWebhookNotifier
	mongo                                                       *irepositories.MongoClient
)

func init() {
	dicordWebhookId = os.Getenv("DISCORD_WEBHOOK_ID")
	discordWebhookToken = os.Getenv("DISCORD_WEBHOOK_TOKEN")
	mongoConnectionString = os.Getenv("MONGO_CONNECTION_STRING")
	ctx := context.TODO()
	mongodbClient, err := irepositories.NewClient(ctx, mongoConnectionString)
	if err != nil {
		log.WithError(err).Fatalln("can't create mongodb client")
	}
	mongo = mongodbClient
	redditParser := iparsers.NewRedditParser([]string{"dotnet", "csharp", "fsharp", "golang", "python", "node", "javascript", "devops"})
	hackernewsParser := iparsers.NewHackerNewsArticleParser()
	dotnetomaniakParser := iparsers.NewDotnetoManiakParser()
	parsers := []parsers.ArticlesParser{redditParser, hackernewsParser, dotnetomaniakParser}
	repo := irepositories.NewMongoArticlesRepository(mongodbClient, "Articles")
	articlesProvider := providers.NewArticlesProvider(parsers)
	discord, err := notifiers.NewDiscordWebhookNotifier(dicordWebhookId, discordWebhookToken)
	if err != nil {
		log.WithError(err).Fatalln("error creating discord notifier")
	}
	discordWebhookNotifier = discord
	bradcaster := notifications.NewBroadcaster([]notifications.Notifier{discordWebhookNotifier})
	usecase := usecase.NewParseArticlesAndSendItUseCase(articlesProvider, repo, bradcaster)
	parser = handler.NewParseAndSendArticlesHandler(usecase)
}

func main() {
	ctx := context.TODO()
	defer mongo.Close(ctx)
	defer discordWebhookNotifier.Close()
	lambda.Start(parser.Handle)
}
