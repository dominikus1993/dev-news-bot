package usecase

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
)

type ParseArticlesAndSendItUseCase struct {
	articlesProvider providers.ArticlesProvider
	repository       repositories.ArticlesWriter
	broadcaster      notifications.Broadcaster
}

func NewParseArticlesAndSendItUseCase(articlesProvider providers.ArticlesProvider, repository repositories.ArticlesWriter, broadcaster notifications.Broadcaster) *ParseArticlesAndSendItUseCase {
	return &ParseArticlesAndSendItUseCase{articlesProvider: articlesProvider, repository: repository, broadcaster: broadcaster}
}

func (u *ParseArticlesAndSendItUseCase) Execute(ctx context.Context, articlesQuantity int) error {
	articles := u.articlesProvider.Provide(ctx)
	randomArticles := model.TakeRandomArticles(ctx, articles, articlesQuantity)
	err := u.repository.Save(ctx, randomArticles)
	if err != nil {
		return err
	}
	return u.broadcaster.Broadcast(ctx, randomArticles)
}
