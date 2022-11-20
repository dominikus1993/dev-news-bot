package providers

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/filters"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
)

type ArticlesProvider interface {
	Provide(ctx context.Context) model.ArticlesStream
}

type articlesProvider struct {
	parsers    []parsers.ArticlesParser
	filter     filters.ArticlesFilter
	repository repositories.ArticlesReader
}

func NewArticlesProvider(repository repositories.ArticlesReader, filter filters.ArticlesFilter, parsers ...parsers.ArticlesParser) *articlesProvider {
	return &articlesProvider{parsers: parsers, repository: repository, filter: filter}
}

func (u *articlesProvider) filterValid(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return channels.Filter(ctx, articles, func(ctx context.Context, article model.Article) bool {
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
	filteredArticles := f.filter.Where(ctx, uniqueArticles)
	newArticles := f.repository.FilterNew(ctx, filteredArticles)
	return newArticles
}
