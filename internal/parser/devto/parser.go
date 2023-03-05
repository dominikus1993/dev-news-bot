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

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"

type devtoresponse []struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func getTagUrl(tag string) string {
	return fmt.Sprintf("https://dev.to/api/articles?tag=%s", tag)
}

func parseTag(ctx context.Context, client *http.Client, tag string) (*devtoresponse, error) {
	url := getTagUrl(tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
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
	client *http.Client
	tags   []string
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func NewDevToParser(tags []string) *devtoParser {
	return &devtoParser{
		client: getClient(),
		tags:   tags,
	}
}

func mapPostToArticle(sub *devtoresponse) []model.Article {
	tag := *sub
	articles := make([]model.Article, len(tag))
	for i, post := range tag {
		articles[i] = model.NewArticleWithContent(post.Title, post.URL, post.Description)
	}
	return articles
}

func (p *devtoParser) parseAll(ctx context.Context) chan []model.Article {
	stream := make(chan []model.Article, 10)
	go func() {
		var wg sync.WaitGroup
		for _, sub := range p.tags {
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
	}()
	return stream
}

func flatten(stream chan []model.Article) model.ArticlesStream {
	result := make(chan model.Article, 100)
	go func() {
		defer close(result)
		var wg sync.WaitGroup
		for articles := range stream {
			wg.Add(1)
			go func(articless []model.Article, stream chan model.Article, wait *sync.WaitGroup) {
				defer wait.Done()
				for _, article := range articless {
					stream <- article
				}
			}(articles, result, &wg)
		}
		wg.Wait()
	}()
	return result
}

func (p *devtoParser) Parse(ctx context.Context) model.ArticlesStream {
	return flatten(p.parseAll(ctx))
}
