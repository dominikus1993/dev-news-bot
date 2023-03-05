package echojs

import (
	"context"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
const url = "https://www.echojs.com/"
const source = "echo.js"

type echojsParser struct {
}

func NewEechoJsParser() *echojsParser {
	return &echojsParser{}
}

func (parser *echojsParser) Parse(ctx context.Context) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		c := colly.NewCollector(colly.Async(true), colly.UserAgent(userAgent))
		c.OnHTML("article h2 a", func(e *colly.HTMLElement) {
			title := e.Text
			link := e.Attr("href")
			result <- model.NewArticle(title, link, source)
		})
		c.SetRequestTimeout(time.Second * 30)
		c.UserAgent = "devnews-bot"
		err := c.Visit(url)
		if err != nil {
			log.WithError(err).Errorln("Error while parsing dotnetomaniak")
		}
		c.Wait()
		close(result)
	}()
	return result
}
