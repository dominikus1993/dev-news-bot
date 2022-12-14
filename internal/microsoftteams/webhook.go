package microsoftteams

import (
	"context"
	"fmt"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
)

type TeamsWebhookNotifier struct {
	webhookUrl string
	client     *goteamsnotify.TeamsClient
}

func NewDiscordWebhookNotifier(webhookUrl string) (*TeamsWebhookNotifier, error) {
	client := goteamsnotify.NewTeamsClient()
	return &TeamsWebhookNotifier{
		webhookUrl: webhookUrl,
		client:     client,
	}, nil
}

func createDiscordEmbedsFromArticles(articles []model.Article) *messagecard.MessageCard {
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Witam serdecznie"
	msgCard.Text = "Oto nowe newsy"
	for _, article := range articles {
		pa, _ := messagecard.NewPotentialAction(
			messagecard.PotentialActionOpenURIType,
			article.GetTitle(),
		)
		pa.PotentialActionOpenURI.Targets =
			[]messagecard.PotentialActionOpenURITarget{
				{
					OS:  "default",
					URI: article.GetLink(),
				},
			}
		section := messagecard.Section{Title: article.GetTitle(), Text: article.GetContent()}
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
