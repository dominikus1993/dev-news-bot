package parser

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/gocolly/colly/v2"
)

const (
	hackerNewsURL string = "https://news.ycombinator.com"
)

type hackerNewsArticleParser struct {
}

func NewHackerNewsArticleParser() *hackerNewsArticleParser {
	return &hackerNewsArticleParser{}
}

func (p *hackerNewsArticleParser) Parse(ctx context.Context) ([]model.Article, error) {
	result := make([]model.Article, 0)
	c := colly.NewCollector()
	c.OnHTML(".titlelink", func(e *colly.HTMLElement) {
		result = append(result, model.NewArticle(e.Text, e.Attr("href")))
	})
	err := c.Visit(hackerNewsURL)
	if err != nil {
		return nil, err
	}
	return result, nil
}
