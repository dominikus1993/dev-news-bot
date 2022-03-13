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
	dotnetomaniakNewsURL    string = "dotnetomaniak.pl"
	dotnetomaniakNewsScheme string = "https"
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
		c := colly.NewCollector()
		c.OnHTML(".article", func(e *colly.HTMLElement) {
			title := e.ChildText(".title .taggedlink span")
			link := geDotnetomaniakLink(e.ChildAttr(".title .taggedlink", "href"))
			content := e.ChildText(".description p span")
			result <- model.NewArticleWithContent(title, link, content)
		})
		c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
			log.Debugf("Redirecting to %s", req.URL.String())
			return nil
		})
		url := fmt.Sprintf("%s://%s/", dotnetomaniakNewsScheme, dotnetomaniakNewsURL)
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
