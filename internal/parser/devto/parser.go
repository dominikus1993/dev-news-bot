package devto

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dominikus1993/dev-news-bot/internal/parser/utils"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
const source = "dev.to"

type devtoresponse []struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func getTagUrl(tag string) string {
	return fmt.Sprintf("https://dev.to/api/articles?tag=%s", tag)
}

func parseTag(ctx context.Context, tag string) (*devtoresponse, error) {
	url := getTagUrl(tag)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent(userAgent)
	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error while parsing hackernews top articless, status: %d", resp.StatusCode())
	}
	body := resp.Body()
	var sub devtoresponse
	if err := json.Unmarshal(body, &sub); err != nil {
		return nil, err
	}
	l := len(sub)
	slog.InfoContext(ctx, "Parsed posts from tag", slog.String("tag", tag), slog.Int("count", l))
	return &sub, nil
}

type devtoParser struct {
	tags []string
}

func NewDevToParser(tags []string) *devtoParser {
	return &devtoParser{
		tags: tags,
	}
}

func streamArticles(sub *devtoresponse, stream chan<- model.Article) {
	tag := *sub
	for _, post := range tag {
		stream <- model.NewArticleWithContent(post.Title, post.URL, post.Description, source)
	}
}

func (p *devtoParser) parseArticles(ctx context.Context, stream chan<- model.Article) {
	var wg sync.WaitGroup
	for _, sub := range p.tags {
		s := sub
		wg.Go(func() {
			res, err := parseTag(ctx, s)
			if err != nil {
				slog.ErrorContext(ctx, "Error while parsing tag", slog.String("tag", s), slog.Any("error", err))
			} else {
				streamArticles(res, stream)
			}
		})
	}
	wg.Wait()
}

func (p *devtoParser) Parse(ctx context.Context) model.ArticlesStream {
	return utils.Parse(ctx, p.parseArticles)
}
