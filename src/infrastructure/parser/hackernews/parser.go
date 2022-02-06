package hackernews

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const (
	hackerNewsURL    string = "news.ycombinator.com"
	hackerNewsScheme string = "https"
)

type hackerNewsArticleParser struct {
}

func NewHackerNewsArticleParser() *hackerNewsArticleParser {
	return &hackerNewsArticleParser{}
}

func getLink(e *colly.HTMLElement) string {
	link := e.Attr("href")
	url, _ := url.Parse(link)
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	if url.Host == "" {
		url.Host = hackerNewsURL
	}
	return url.String()
}

func (p *hackerNewsArticleParser) Parse(ctx context.Context) ([]model.Article, error) {
	result := make([]model.Article, 0)
	c := colly.NewCollector()
	c.OnHTML(".titlelink", func(e *colly.HTMLElement) {
		result = append(result, model.NewArticle(e.Text, getLink(e)))
	})
	err := c.Visit(fmt.Sprintf("%s://%s", hackerNewsScheme, hackerNewsURL))
	if err != nil {
		return nil, err
	}
	log.WithContext(ctx).WithField("quantity", len(result)).Infoln("Parsed hackernews articles")
	return result, nil
}
