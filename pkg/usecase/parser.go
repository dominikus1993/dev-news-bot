package usecase

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	log "github.com/sirupsen/logrus"
)

type ParseArticlesAndSendItUseCase struct {
	articlesProvider providers.ArticlesProvider
	repository       repositories.ArticleRepository
	broadcaster      notifications.Broadcaster
}

func NewParseArticlesAndSendItUseCase(articlesProvider providers.ArticlesProvider, repository repositories.ArticleRepository, broadcaster notifications.Broadcaster) *ParseArticlesAndSendItUseCase {
	return &ParseArticlesAndSendItUseCase{articlesProvider: articlesProvider, repository: repository, broadcaster: broadcaster}
}

func (u *ParseArticlesAndSendItUseCase) filterNewArticles(ctx context.Context, articles []model.Article) []model.Article {
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

func (u *ParseArticlesAndSendItUseCase) filterValid(ctx context.Context, articles []model.Article) []model.Article {
	filteredArticles := make([]model.Article, 0, len(articles))
	for _, article := range articles {
		if article.IsValid() {
			filteredArticles = append(filteredArticles, article)
		}
	}
	return filteredArticles
}

func (u *ParseArticlesAndSendItUseCase) Execute(ctx context.Context, articlesQuantity int) error {
	articles, err := u.articlesProvider.Provide(ctx)
	if err != nil {
		return err
	}
	log.Infoln("Found articles:", len(articles))
	validArticles := u.filterValid(ctx, articles)
	log.Infoln("Found valid articles:", len(validArticles))
	newArticles := u.filterNewArticles(ctx, validArticles)
	log.Infoln("Found new articles:", len(newArticles))
	randomArticles := model.TakeRandomArticles(newArticles, articlesQuantity)
	err = u.broadcaster.Broadcast(ctx, randomArticles)
	if err != nil {
		return err
	}
	return u.repository.Save(ctx, randomArticles)
}
