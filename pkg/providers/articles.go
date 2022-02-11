package providers

import (
	"context"
	"sync"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

type ArticlesProvider interface {
	Provide(ctx context.Context) ([]model.Article, error)
}

type articlesProvider struct {
	parsers []parsers.ArticlesParser
}

func NewArticlesProvider(parsers []parsers.ArticlesParser) *articlesProvider {
	return &articlesProvider{parsers: parsers}
}

func (f *articlesProvider) parse(ctx context.Context, parser parsers.ArticlesParser, stream chan<- []model.Article, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := parser.Parse(ctx)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("Error while parsing articles")
	} else {
		stream <- res
	}
}

func (f *articlesProvider) parseAll(ctx context.Context, stream chan<- []model.Article) {
	var wg sync.WaitGroup
	for _, parser := range f.parsers {
		wg.Add(1)
		go f.parse(ctx, parser, stream, &wg)
	}
	wg.Wait()
	close(stream)
}

func (f *articlesProvider) Provide(ctx context.Context) ([]model.Article, error) {
	stream := make(chan []model.Article, 10)
	result := make([]model.Article, 0)
	go f.parseAll(ctx, stream)
	for articles := range stream {
		result = append(result, articles...)
	}
	return result, nil
}
