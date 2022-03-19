package providers

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	log "github.com/sirupsen/logrus"
)

type ArticlesProvider interface {
	Provide(ctx context.Context) model.ArticlesStream
}

type articlesProvider struct {
	parsers    []parsers.ArticlesParser
	repository repositories.ArticlesReader
}

func NewArticlesProvider(parsers []parsers.ArticlesParser, repository repositories.ArticlesReader) *articlesProvider {
	return &articlesProvider{parsers: parsers, repository: repository}
}

func (u *articlesProvider) filterNewArticles(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return channels.Filter(ctx, articles, func(ctx context.Context, article *model.Article) bool {
		isNew, err := u.repository.IsNew(ctx, article)
		if err != nil {
			log.WithField("ArticleLink", article.Link).WithError(err).WithContext(ctx).Error("error while checking if article exists")
			return false
		}
		return isNew
	})
}

func (u *articlesProvider) filterValid(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return channels.Filter(ctx, articles, func(ctx context.Context, article *model.Article) bool {
		return article.IsValid()
	})
}

func (f *articlesProvider) Provide(ctx context.Context) model.ArticlesStream {
	streams := make([]<-chan model.Article, 0, len(f.parsers))
	for _, parser := range f.parsers {
		streams = append(streams, parser.Parse(ctx))
	}
	articles := channels.FanIn(ctx, streams...)
	validArticles := f.filterValid(ctx, articles)
	uniqueArticles := model.UniqueArticles(validArticles)
	newArticles := f.filterNewArticles(ctx, uniqueArticles)
	return newArticles
}
