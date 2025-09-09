package cmd

import (
	"context"
	"log/slog"

	"github.com/dominikus1993/dev-news-bot/internal/console"
	"github.com/dominikus1993/dev-news-bot/internal/discord"
	"github.com/dominikus1993/dev-news-bot/internal/language"
	"github.com/dominikus1993/dev-news-bot/internal/mongo"
	"github.com/dominikus1993/dev-news-bot/internal/parser/devto"
	"github.com/dominikus1993/dev-news-bot/internal/parser/dotnetomaniak"
	"github.com/dominikus1993/dev-news-bot/internal/parser/echojs"
	"github.com/dominikus1993/dev-news-bot/internal/parser/hackernews"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/usecase"
	"github.com/urfave/cli/v3"
)

type notifiers struct {
	discord *discord.DiscordWebhookNotifier
	pprint  *console.PPrintNotifier
}

func createNotifiers(cmd *ParseArgs) (notifiers, error) {
	var notifiers notifiers
	if cmd.dicordWebhookId != "" && cmd.discordWebhookToken != "" {
		discordn, err := discord.NewDiscordWebhookNotifier(cmd.dicordWebhookId, cmd.discordWebhookToken)
		if err != nil {
			return notifiers, err
		}
		notifiers.discord = discordn
	}
	notifiers.pprint = console.NewPPrintNotifier()
	return notifiers, nil
}

func (n notifiers) createBroadcaster() notifications.Broadcaster {
	var notifiers []notifications.Notifier
	if n.discord != nil {
		notifiers = append(notifiers, n.discord)
	}
	if n.pprint != nil {
		notifiers = append(notifiers, n.pprint)
	}
	return notifications.NewBroadcaster(notifiers...)
}

func (n notifiers) close() {
	if n.discord != nil {
		n.discord.Close()
	}
}

type ParseArgs struct {
	quantity              int
	dicordWebhookId       string
	discordWebhookToken   string
	mongoConnectionString string
	teamsWebhookUrl       string
}

func NewParseArgs(context *cli.Command) *ParseArgs {
	dicordWebhookId := context.String("dicord-webhook-id")
	discordWebhhokToken := context.String("discord-webhook-token")
	quantity := context.Int("quantity")
	mongo := context.String("mongo-connection-string")
	teams := context.String("teams-webhook-url")
	return &ParseArgs{dicordWebhookId: dicordWebhookId, discordWebhookToken: discordWebhhokToken, quantity: int(quantity), mongoConnectionString: mongo, teamsWebhookUrl: teams}
}

func Parse(ctx context.Context, cmd *cli.Command) error {
	p := NewParseArgs(cmd)
	slog.InfoContext(ctx, "Parse Articles And Send It")
	mongodbClient, err := mongo.NewClient(ctx, p.mongoConnectionString, "Articles")
	if err != nil {
		slog.ErrorContext(ctx, "can't create mongodb client", "error", err)
		return cli.Exit("can't create mongodb client", 1)
	}
	defer mongodbClient.Close(ctx)
	devtoParser := devto.NewDevToParser([]string{"dotnet", "csharp", "fsharp", "golang", "python", "node", "javascript", "devops", "rust", "aws", "vlang", "typescript", "react"})
	hackernewsParser := hackernews.NewHackerNewsArticleParser(50)
	dotnetomaniakParser := dotnetomaniak.NewDotnetoManiakParser()
	echojsp := echojs.NewEechoJsParser()
	repo := mongo.NewMongoArticlesRepository(mongodbClient)
	languageFilter := language.NewLanguageFilter()
	articlesProvider := providers.NewArticlesProvider(repo, languageFilter, hackernewsParser, dotnetomaniakParser, devtoParser, echojsp)
	notifiers, err := createNotifiers(p)
	if err != nil {
		slog.ErrorContext(ctx, "can't create notifiers", "error", err)
		return cli.Exit("can't create notifiers", 1)
	}
	defer notifiers.close()

	bradcaster := notifiers.createBroadcaster()
	usecase := usecase.NewParseArticlesAndSendItUseCase(articlesProvider, repo, bradcaster)

	err = usecase.Execute(ctx, p.quantity)

	if err != nil {
		slog.ErrorContext(ctx, "Error while parsing articles", "error", err)
		return cli.Exit("Error while parsing articles", 0)
	}
	slog.InfoContext(ctx, "Parsing articles finished")
	return nil
}
