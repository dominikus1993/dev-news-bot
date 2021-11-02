package parser

import (
	"context"
	"fmt"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/gocolly/colly/v2"
)

const (
	dotnetomaniakNewsURL    string = "dotnetomaniak.pl"
	dotnetomaniakNewsScheme string = "https"
)

type dotnetoManiakParser struct {
}

func NewDotnetoManiakParser() *dotnetoManiakParser {
	return &dotnetoManiakParser{}
}

func (p *dotnetoManiakParser) Parse(ctx context.Context) ([]model.Article, error) {
	result := make([]model.Article, 0)
	c := colly.NewCollector()
	c.OnHTML(".article", func(e *colly.HTMLElement) {
		result = append(result, model.NewArticle(e.Text, getLink(e)))
	})
	err := c.Visit(fmt.Sprintf("%s://%s/", dotnetomaniakNewsScheme, dotnetomaniakNewsURL))
	if err != nil {
		return nil, err
	}
	return result, nil
}
