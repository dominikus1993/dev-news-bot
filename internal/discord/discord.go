package discord

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

const discordMessageSizeLimit = 10

type DiscordWebhookNotifier struct {
	webhookID    string
	webhookToken string
	client       *discordgo.Session
}

func NewDiscordWebhookNotifier(webhookID, webhookToken string) (*DiscordWebhookNotifier, error) {
	session, err := discordgo.New("")
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
		slog.Error("Error while closing discord session", slog.Any("error", err))
	}
}

func createDiscordEmbedsFromArticles(articles []model.Article) []*discordgo.MessageEmbed {
	embeds := make([]*discordgo.MessageEmbed, 0)
	for _, article := range articles {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       article.GetTitle(),
			Description: article.GetContent(),
			URL:         article.GetLink(),
			Color:       0x00ff00,
		})
	}
	return embeds
}

func (not *DiscordWebhookNotifier) Notify(ctx context.Context, articles []model.Article) error {
	var result error
	log.Infoln("Notifying via discord")
	embeds := createDiscordEmbedsFromArticles(articles)
	if len(embeds) <= discordMessageSizeLimit {
		return not.send(ctx, embeds)
	}

	chunks := lo.Chunk(embeds, 10)
	for _, chunk := range chunks {
		err := not.send(ctx, chunk)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func (not *DiscordWebhookNotifier) send(_ context.Context, embeds []*discordgo.MessageEmbed) error {
	msg := discordgo.WebhookParams{Content: "Witam serdecznie, oto nowe newsy", Embeds: embeds}
	_, err := not.client.WebhookExecute(not.webhookID, not.webhookToken, true, &msg)
	if err != nil {
		return fmt.Errorf("error while sending webhook: %w", err)
	}
	return nil
}
