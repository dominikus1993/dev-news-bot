package devto

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

type devtoresponse []struct {
	TypeOf                 string      `json:"type_of"`
	ID                     int         `json:"id"`
	Title                  string      `json:"title"`
	Description            string      `json:"description"`
	CoverImage             string      `json:"cover_image"`
	ReadablePublishDate    string      `json:"readable_publish_date"`
	SocialImage            string      `json:"social_image"`
	TagList                []string    `json:"tag_list"`
	Tags                   string      `json:"tags"`
	Slug                   string      `json:"slug"`
	Path                   string      `json:"path"`
	URL                    string      `json:"url"`
	CanonicalURL           string      `json:"canonical_url"`
	CommentsCount          int         `json:"comments_count"`
	PositiveReactionsCount int         `json:"positive_reactions_count"`
	PublicReactionsCount   int         `json:"public_reactions_count"`
	CollectionID           interface{} `json:"collection_id"`
	CreatedAt              time.Time   `json:"created_at"`
	EditedAt               time.Time   `json:"edited_at"`
	CrosspostedAt          interface{} `json:"crossposted_at"`
	PublishedAt            time.Time   `json:"published_at"`
	LastCommentAt          time.Time   `json:"last_comment_at"`
	PublishedTimestamp     time.Time   `json:"published_timestamp"`
	ReadingTimeMinutes     int         `json:"reading_time_minutes"`
	User                   struct {
		Name            string `json:"name"`
		Username        string `json:"username"`
		TwitterUsername string `json:"twitter_username"`
		GithubUsername  string `json:"github_username"`
		WebsiteURL      string `json:"website_url"`
		ProfileImage    string `json:"profile_image"`
		ProfileImage90  string `json:"profile_image_90"`
	} `json:"user"`
	Organization struct {
		Name           string `json:"name"`
		Username       string `json:"username"`
		Slug           string `json:"slug"`
		ProfileImage   string `json:"profile_image"`
		ProfileImage90 string `json:"profile_image_90"`
	} `json:"organization"`
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
	for s := range stream {
		articles = append(articles, s...)
	}
	return articles, nil
}
