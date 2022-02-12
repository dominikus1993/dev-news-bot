package devto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
)

type devtoresponse []struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func parseTag(ctx context.Context, client *http.Client, tag string) (*devtoresponse, error) {
	url := fmt.Sprintf("https://dev.to/api/articles?tag=%s", tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "dev-news-bot")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while parsing tag: %s, status: %s", tag, resp.Status)
	}

	defer resp.Body.Close()
	var sub devtoresponse
	if err := json.NewDecoder(resp.Body).Decode(&sub); err != nil {
		return nil, err
	}
	l := len(sub)
	log.WithContext(ctx).Infof("Parsed %d posts from tag %s", l, tag)
	return &sub, nil
}

type devtoParser struct {
	client     *http.Client
	subreddits []string
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func NewDevToParser(subreddits []string) *devtoParser {
	return &devtoParser{
		client:     getClient(),
		subreddits: subreddits,
	}
}

func mapPostToArticle(sub *devtoresponse) []model.Article {
	articles := make([]model.Article, len(*sub))
	for i, post := range *sub {
		articles[i] = model.NewArticleWithContent(post.Title, post.URL, post.Description)
	}
	return articles
}

func (p *devtoParser) parseAll(ctx context.Context, stream chan []model.Article) {
	var wg sync.WaitGroup
	for _, sub := range p.subreddits {
		wg.Add(1)
		go func(s string, wait *sync.WaitGroup) {
			defer wg.Done()
			res, err := parseTag(ctx, p.client, s)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Error while parsing tag: %s", s)
			} else {
				stream <- mapPostToArticle(res)
			}
		}(sub, &wg)
	}
	wg.Wait()
	close(stream)
}

func (p *devtoParser) Parse(ctx context.Context) ([]model.Article, error) {
	stream := make(chan []model.Article)
	go p.parseAll(ctx, stream)
	articles := make([]model.Article, 0)
	for {
		select {
		case <-ctx.Done():
			return articles, nil
		case v, ok := <-stream:
			if !ok {
				return articles, nil
			}
			articles = append(articles, v...)
		}
	}
}
