package dotnetomaniak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const (
	userAgent               string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
	dotnetomaniakNewsURL    string = "dotnetomaniak.pl"
	dotnetomaniakNewsScheme string = "https"
	source                  string = "dotnetomaniak"
)

type dotnetoManiakParser struct {
}

func NewDotnetoManiakParser() *dotnetoManiakParser {
	return &dotnetoManiakParser{}
}

func geDotnetomaniakLink(link string) string {
	url, _ := url.Parse(link)
	if url.Scheme == "" {
		url.Scheme = dotnetomaniakNewsScheme
	}
	if url.Host == "" {
		url.Host = dotnetomaniakNewsURL
	}
	return url.String()
}

func (p *dotnetoManiakParser) Parse(ctx context.Context) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		defer close(result)
		c := colly.NewCollector(colly.Async(true), colly.UserAgent(userAgent))
		c.OnHTML(".article", func(e *colly.HTMLElement) {
			title := e.ChildText(".title .taggedlink span")
			link := geDotnetomaniakLink(e.ChildAttr(".title .taggedlink", "href"))
			content := e.ChildText(".description p span")
			result <- model.NewArticleWithContent(title, link, content, source)
		})
		c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
			log.Debugf("Redirecting to %s", req.URL.String())
			return nil
		})
		c.OnError(func(r *colly.Response, err error) {
			log.WithError(err).Errorln("can't parse dotnetomaniak")
		})
		url := fmt.Sprintf("%s://%s/", dotnetomaniakNewsScheme, dotnetomaniakNewsURL)
		c.SetRequestTimeout(time.Second * 30)
		c.UserAgent = "devnews-bot"
		err := c.Visit(url)
		if err != nil {
			log.WithError(err).Errorln("error while parsing dotnetomaniak")
		}
		c.Wait()

	}()
	return result
}
