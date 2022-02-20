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

func pipe(ctx context.Context, articles model.ArticlesStream, f func(ctx context.Context, article model.Article, articles chan<- model.Article)) model.ArticlesStream {
	filteredArticles := make(chan model.Article, 100)
	go func() {
		for article := range articles {
			f(ctx, article, filteredArticles)
		}
		close(filteredArticles)
	}()
	return filteredArticles
}

func (u *ParseArticlesAndSendItUseCase) filterNewArticles(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return pipe(ctx, articles, func(ctx context.Context, article model.Article, articles chan<- model.Article) {
		isNew, err := u.repository.IsNew(ctx, article)
		if err != nil {
			log.WithField("ArticleLink", article.Link).WithError(err).WithContext(ctx).Error("error while checking if article exists")
		}
		if isNew {
			articles <- article
		}
	})
}

func (u *ParseArticlesAndSendItUseCase) filterValid(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return pipe(ctx, articles, func(ctx context.Context, article model.Article, articles chan<- model.Article) {
		if article.IsValid() {
			articles <- article
		}
	})
}

func (u *ParseArticlesAndSendItUseCase) Execute(ctx context.Context, articlesQuantity int) error {
	articles := u.articlesProvider.Provide(ctx)
	validArticles := u.filterValid(ctx, articles)
	newArticles := u.filterNewArticles(ctx, validArticles)
	randomArticles := model.TakeRandomArticles(newArticles, articlesQuantity)
	err := u.broadcaster.Broadcast(ctx, randomArticles)
	if err != nil {
		return err
	}
	return u.repository.Save(ctx, randomArticles)
}
