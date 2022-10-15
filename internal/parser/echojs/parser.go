package echojs

import (
	"context"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const url = "https://www.echojs.com/"

type echojsParser struct {
}

func NewEechoJsParser() *echojsParser {
	return &echojsParser{}
}

func (parser *echojsParser) Parse(ctx context.Context) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		c := colly.NewCollector()
		c.OnHTML("article h2 a", func(e *colly.HTMLElement) {
			title := e.Text
			link := e.Attr("href")
			result <- model.NewArticle(title, link)
		})
		c.SetRequestTimeout(time.Second * 30)
		c.UserAgent = "devnews-bot"
		err := c.Visit(url)
		if err != nil {
			log.WithError(err).Errorln("Error while parsing dotnetomaniak")
		}
		close(result)
	}()
	return result
}
