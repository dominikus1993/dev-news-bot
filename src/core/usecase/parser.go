package usecase

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/dominikus1993/dev-news-bot/src/core/notifications"
	"github.com/dominikus1993/dev-news-bot/src/core/providers"
	"github.com/dominikus1993/dev-news-bot/src/core/repositories"
	log "github.com/sirupsen/logrus"
)

type parseArticlesAndSendItUseCase struct {
	articlesProvider providers.ArticlesProvider
	repository       repositories.IArticleRepository
	broadcaster      notifications.Broadcaster
}

func NewParseArticlesAndSendItUseCase(articlesProvider providers.ArticlesProvider, repository repositories.IArticleRepository, broadcaster notifications.Broadcaster) *parseArticlesAndSendItUseCase {
	return &parseArticlesAndSendItUseCase{articlesProvider: articlesProvider, repository: repository, broadcaster: broadcaster}
}

func (u *parseArticlesAndSendItUseCase) filterNewArticles(ctx context.Context, articles []model.Article) []model.Article {
	filteredArticles := make([]model.Article, 0, len(articles))
	for _, article := range articles {
		isNew, err := u.repository.IsNew(ctx, article)
		if err != nil {
			log.WithField("ArticleLink", article.Link).WithError(err).WithContext(ctx).Error("error while checking if article exists")
		}
		if isNew {
			filteredArticles = append(filteredArticles, article)
		}
	}
	return filteredArticles
}

func (u *parseArticlesAndSendItUseCase) filterValid(ctx context.Context, articles []model.Article) []model.Article {
	filteredArticles := make([]model.Article, 0, len(articles))
	for _, article := range articles {
		if article.IsValid() {
			filteredArticles = append(filteredArticles, article)
		}
	}
	return filteredArticles
}

func (u *parseArticlesAndSendItUseCase) Execute(ctx context.Context, articlesQuantity int) error {
	articles, err := u.articlesProvider.Provide(ctx)
	if err != nil {
		return err
	}
	validArticles := u.filterValid(ctx, articles)
	newArticles := u.filterNewArticles(ctx, validArticles)
	randomArticles := model.TakeRandomArticles(newArticles, articlesQuantity)
	err = u.repository.Save(ctx, randomArticles)
	if err != nil {
		return err
	}
	return u.broadcaster.Broadcast(ctx, randomArticles)
}
