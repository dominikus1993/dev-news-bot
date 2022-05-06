package microsoftteams

import (
	"context"
	"fmt"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
)

type TeamsWebhookNotifier struct {
	webhookUrl string
	client     goteamsnotify.API
}

func NewDiscordWebhookNotifier(webhookUrl string) (*TeamsWebhookNotifier, error) {
	client := goteamsnotify.NewClient()
	return &TeamsWebhookNotifier{
		webhookUrl: webhookUrl,
		client:     client,
	}, nil
}

func createDiscordEmbedsFromArticles(articles []model.Article) goteamsnotify.MessageCard {
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = "Witam serdecznie, oto nowe newsy"
	msgCard.Text = "Oto nowe newsy"
	for _, article := range articles {
		pa, _ := goteamsnotify.NewMessageCardPotentialAction(
			goteamsnotify.PotentialActionOpenURIType,
			article.GetTitle(),
		)
		pa.MessageCardPotentialActionOpenURI.Targets =
			[]goteamsnotify.MessageCardPotentialActionOpenURITarget{
				{
					OS:  "default",
					URI: article.GetLink(),
				},
			}
		section := goteamsnotify.MessageCardSection{Title: article.GetTitle(), Text: article.GetContent()}
		section.AddPotentialAction(pa)
		msgCard.AddSection(&section)
	}
	return msgCard
}

func (not *TeamsWebhookNotifier) Notify(ctx context.Context, articles []model.Article) error {
	log.Infoln("notifying via teams")
	msg := createDiscordEmbedsFromArticles(articles)
	err := not.client.Send(not.webhookUrl, msg)
	if err != nil {
		return fmt.Errorf("error while sending teams webhook: %w", err)
	}
	return nil
}
