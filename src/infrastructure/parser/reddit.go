package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	log "github.com/sirupsen/logrus"
)

type subreddit struct {
	Data struct {
		Children []struct {
			Data struct {
				Title   string `json:"title"`
				URL     string `json:"url"`
				Content string `json:"selftext"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func parseSubreddit(ctx context.Context, client *http.Client, subr string) (*subreddit, error) {
	resp, err := client.Get(fmt.Sprintf("https://www.reddit.com/r/%s.json?limit=10", subr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var sub subreddit
	if err := json.NewDecoder(resp.Body).Decode(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

type redditParser struct {
	client     *http.Client
	subreddits []string
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func NewRedditParser(subreddits []string) *redditParser {
	return &redditParser{
		client:     getClient(),
		subreddits: subreddits,
	}
}

func mapPostToArticle(sub *subreddit) []model.Article {
	articles := make([]model.Article, len(sub.Data.Children))
	for i, post := range sub.Data.Children {
		articles[i] = model.NewArticleWithContent(post.Data.Title, post.Data.URL, post.Data.Content)
	}
	return articles
}

func (p *redditParser) parseAll(ctx context.Context, stream chan []model.Article) {
	var wg sync.WaitGroup
	for _, sub := range p.subreddits {
		wg.Add(1)
		go func(sub string, wait *sync.WaitGroup) {
			defer wg.Done()
			res, err := parseSubreddit(ctx, p.client, sub)
			if err != nil {
				log.WithError(err).Errorf("Error while parsing subreddit: %s", sub)
			}
			stream <- mapPostToArticle(res)
		}(sub, &wg)
	}
	wg.Wait()
	close(stream)
}

func (p *redditParser) Parse(ctx context.Context) ([]model.Article, error) {
	stream := make(chan []model.Article)
	go p.parseAll(ctx, stream)
	var articles []model.Article
	for s := range stream {
		articles = append(articles, s...)
	}
	return articles, nil
}
