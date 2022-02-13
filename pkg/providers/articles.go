package providers

import (
	"context"
	"sync"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

type ArticlesProvider interface {
	Provide(ctx context.Context) model.ArticlesStream
}

type articlesProvider struct {
	parsers []parsers.ArticlesParser
}

func NewArticlesProvider(parsers []parsers.ArticlesParser) *articlesProvider {
	return &articlesProvider{parsers: parsers}
}

func fanIn(ctx context.Context, stream ...chan model.Article) chan model.Article {
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

func (f *articlesProvider) parse(ctx context.Context, parser parsers.ArticlesParser) chan model.Article {
	stream := make(chan model.Article, 10)
	go func() {
		res, err := parser.Parse(ctx)
		if err != nil {
			log.WithError(err).WithContext(ctx).Error("Error while parsing articles")
		} else {
			for _, v := range res {
				select {
				case <-ctx.Done():
					return
				case stream <- v:
				}
			}
		}
		close(stream)
	}()
	return stream
}

func (f *articlesProvider) Provide(ctx context.Context) model.ArticlesStream {
	streams := make([]chan model.Article, 0, len(f.parsers))
	for _, parser := range f.parsers {
		streams = append(streams, f.parse(ctx, parser))
	}
	return fanIn(ctx, streams...)
}
