package echojs

import (
	"context"
	"time"

	"github.com/dominikus1993/dev-news-bot/internal/parser/utils"
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
	return utils.Parse(ctx, parser.parseArticles)
}

func (p *echojsParser) parseArticles(ctx context.Context, result chan<- model.Article) {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent(userAgent))
	c.OnHTML("article h2 a", func(e *colly.HTMLElement) {
		title := e.Text
		link := e.Attr("href")
		article := model.NewArticle(title, link, source)
		if article.IsValid() {
			result <- article
		} else {
			log.WithField("link", article.GetLink()).Warnln("echojs article is not valid")
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.WithError(err).Errorln("can't parse echojs")
	})
	c.SetRequestTimeout(time.Second * 30)
	c.UserAgent = userAgent
	err := c.Visit(url)
	if err != nil {
		log.WithError(err).Errorln("error while parsing echojs")
	}
	c.Wait()
}
