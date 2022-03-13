package providers

import (
	"context"
	"sync"

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

func fanIn(ctx context.Context, stream ...model.ArticlesStream) chan model.Article {
	var wg sync.WaitGroup
	out := make(chan model.Article)
	output := func(c <-chan model.Article) {
		for v := range c {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
		wg.Done()
	}
	wg.Add(len(stream))
	for _, c := range stream {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func filter(ctx context.Context, articles model.ArticlesStream, predicate func(ctx context.Context, article *model.Article) bool) model.ArticlesStream {
	filteredArticles := make(chan model.Article, 100)
	go func() {
		for article := range articles {
			if predicate(ctx, &article) {
				filteredArticles <- article
			}
		}
		close(filteredArticles)
	}()
	return filteredArticles
}

func (u *articlesProvider) filterNewArticles(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return filter(ctx, articles, func(ctx context.Context, article *model.Article) bool {
		isNew, err := u.repository.IsNew(ctx, article)
		if err != nil {
			log.WithField("ArticleLink", article.Link).WithError(err).WithContext(ctx).Error("error while checking if article exists")
			return false
		}
		return isNew
	})
}

func (u *articlesProvider) filterValid(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return filter(ctx, articles, func(ctx context.Context, article *model.Article) bool {
		return article.IsValid()
	})
}

func (f *articlesProvider) Provide(ctx context.Context) model.ArticlesStream {
	streams := make([]model.ArticlesStream, 0, len(f.parsers))
	for _, parser := range f.parsers {
		streams = append(streams, parser.Parse(ctx))
	}
	articles := fanIn(ctx, streams...)
	validArticles := f.filterValid(ctx, articles)
	uniqueArticles := model.UniqueArticles(validArticles)
	newArticles := f.filterNewArticles(ctx, uniqueArticles)
	return newArticles
}
