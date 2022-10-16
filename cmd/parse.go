package cmd

import (
	"context"
	"flag"

	"github.com/dominikus1993/dev-news-bot/internal/console"
	"github.com/dominikus1993/dev-news-bot/internal/discord"
	"github.com/dominikus1993/dev-news-bot/internal/language"
	"github.com/dominikus1993/dev-news-bot/internal/microsoftteams"
	"github.com/dominikus1993/dev-news-bot/internal/mongo"
	"github.com/dominikus1993/dev-news-bot/internal/parser/devto"
	"github.com/dominikus1993/dev-news-bot/internal/parser/dotnetomaniak"
	"github.com/dominikus1993/dev-news-bot/internal/parser/echojs"
	"github.com/dominikus1993/dev-news-bot/internal/parser/hackernews"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/usecase"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type notifiers struct {
	discord *discord.DiscordWebhookNotifier
	pprint  *console.PPrintNotifier
	teams   *microsoftteams.TeamsWebhookNotifier
}

func createNotifiers(cmd *ParseArticlesAndSendIt) (notifiers, error) {
	var notifiers notifiers
	if cmd.dicordWebhookId != "" && cmd.discordWebhookToken != "" {
		discordn, err := discord.NewDiscordWebhookNotifier(cmd.dicordWebhookId, cmd.discordWebhookToken)
		if err != nil {
			return notifiers, err
		}
		notifiers.discord = discordn
	}
	if cmd.teamsWebhookUrl != "" {
		teamsn, err := microsoftteams.NewDiscordWebhookNotifier(cmd.teamsWebhookUrl)
		if err != nil {
			return notifiers, err
		}
		notifiers.teams = teamsn
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
	if n.teams != nil {
		notifiers = append(notifiers, n.teams)
	}
	return notifications.NewBroadcaster(notifiers...)
}

func (n notifiers) close() {
	if n.discord != nil {
		n.discord.Close()
	}
}

type ParseArticlesAndSendIt struct {
	quantity              int
	dicordWebhookId       string
	discordWebhookToken   string
	mongoConnectionString string
	teamsWebhookUrl       string
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
	f.StringVar(&p.teamsWebhookUrl, "teams-webhook-url", "", "teams webhook url")
}

func (p *ParseArticlesAndSendIt) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log.Infoln("Parse Articles And Send It")
	mongodbClient, err := mongo.NewClient(ctx, p.mongoConnectionString, "Articles")
	if err != nil {
		log.WithError(err).Error("can't create mongodb client")
		return subcommands.ExitFailure
	}
	defer mongodbClient.Close(ctx)
	devtoParser := devto.NewDevToParser([]string{"dotnet", "csharp", "fsharp", "golang", "python", "node", "javascript", "devops", "rust", "aws"})
	hackernewsParser := hackernews.NewHackerNewsArticleParser(50)
	dotnetomaniakParser := dotnetomaniak.NewDotnetoManiakParser()
	echojsp := echojs.NewEechoJsParser()
	repo := mongo.NewMongoArticlesRepository(mongodbClient)
	languageFilter := language.NewLanguageFilter()
	articlesProvider := providers.NewArticlesProvider(repo, languageFilter, hackernewsParser, dotnetomaniakParser, devtoParser, echojsp)
	notifiers, err := createNotifiers(p)
	if err != nil {
		log.WithError(err).Error("can't create notifiers")
		return subcommands.ExitFailure
	}
	defer notifiers.close()

	bradcaster := notifiers.createBroadcaster()
	usecase := usecase.NewParseArticlesAndSendItUseCase(articlesProvider, repo, bradcaster)

	err = usecase.Execute(ctx, p.quantity)

	if err != nil {
		log.WithError(err).Error("Error while parsing articles")
		return subcommands.ExitFailure
	}
	log.Infoln("Parsing articles finished")
	return subcommands.ExitSuccess
}
