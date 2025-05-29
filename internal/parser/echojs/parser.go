package echojs

import (
	"context"
	"log/slog"
	"time"

	"github.com/dominikus1993/dev-news-bot/internal/parser/utils"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/gocolly/colly/v2"
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
			slog.Warn("echojs article is not valid", slog.String("link", article.GetLink()))
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("error while parsing echojs", slog.Any("error", err), slog.Int("status_code", r.StatusCode))
	})
	c.SetRequestTimeout(time.Second * 30)
	c.UserAgent = userAgent
	err := c.Visit(url)
	if err != nil {
		slog.Error("error while parsing echojs, stopping parser", slog.Any("error", err))
	}
	c.Wait()
}
