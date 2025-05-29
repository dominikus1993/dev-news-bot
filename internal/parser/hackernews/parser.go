package hackernews

import (
	"context"
	"fmt"

	"github.com/dominikus1993/dev-news-bot/internal/parser/utils"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/go-toolkit/random"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const source = "hackernews"

type hackernewsArticle struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

type hackerNewsArticleParser struct {
	maxArticlesQuantity int
}

func NewHackerNewsArticleParser(maxArticlesQuantity int) *hackerNewsArticleParser {
	return &hackerNewsArticleParser{maxArticlesQuantity: maxArticlesQuantity}
}

func takeRandomArticesIds(ids []int, take int) []int {
	return random.TakeRandomFromSlice(ids, take)
}

func getTopArticlesIds() ([]int, error) {
	const url = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("dev-news-bot")
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
	var res []int
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	l := len(res)
	slog.Info("Parsed ids from hackernews", slog.Int("count", l))
	return res, nil
}

func getArticle(id int) (*hackernewsArticle, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", id)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("dev-news-bot")
	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("error while parsing hackernews article by id: %d, status: %d", id, resp.StatusCode())
	}

	body := resp.Body()
	var res hackernewsArticle
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *hackerNewsArticleParser) Parse(ctx context.Context) model.ArticlesStream {
	return utils.Parse(ctx, p.parseArticles)
}

func (p *hackerNewsArticleParser) parseArticles(ctx context.Context, result chan<- model.Article) {
	ids, err := getTopArticlesIds()
	if err != nil {
		slog.ErrorContext(ctx, "Error while parsing hackernews top articles, stopping parser", slog.Any("error", err))
		return
	}
	ids = takeRandomArticesIds(ids, p.maxArticlesQuantity)
	for _, id := range ids {
		hackerNewsArticle, err := getArticle(id)
		if err != nil {
			slog.ErrorContext(ctx, "Error while parsing hackernews article by id", slog.Int("id", id), slog.Any("error", err))
			continue
		}
		article := model.NewArticle(hackerNewsArticle.Title, hackerNewsArticle.URL, source)
		if article.IsValid() {
			result <- article
		}
	}
}
