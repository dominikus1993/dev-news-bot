package notifiers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
)

type DiscordWebhookNotifier struct {
	webhookID    string
	webhookToken string
	client       *discordgo.Session
}

func NewDiscordWebhookNotifier(webhookID, webhookToken string) (*DiscordWebhookNotifier, error) {
	session, err := discordgo.New()
	if err != nil {
		return nil, err
	}
	return &DiscordWebhookNotifier{
		webhookID:    webhookID,
		webhookToken: webhookToken,
		client:       session,
	}, nil
}

func (not *DiscordWebhookNotifier) Close() {
	err := not.client.Close()
	if err != nil {
		log.WithError(err).Error("Error while closing discord session")
	}
}

func createDiscordEmbedsFromArticles(articles []model.Article) []*discordgo.MessageEmbed {
	embeds := make([]*discordgo.MessageEmbed, 0)
	for _, article := range articles {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       article.Title,
			Description: article.Content,
			URL:         article.Link,
			Color:       0x00ff00,
		})
	}
	return embeds
}

func (not *DiscordWebhookNotifier) Notify(ctx context.Context, articles []model.Article) error {
	msg := discordgo.WebhookParams{Content: "Witam serdecznie, oto nowe newsy", Embeds: createDiscordEmbedsFromArticles(articles)}
	_, err := not.client.WebhookExecute(not.webhookID, not.webhookToken, true, &msg)
	return fmt.Errorf("error while sending webhook: %w", err)
}
