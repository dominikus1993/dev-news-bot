package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
)

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
	client *http.Client
}

func NewHackerNewsArticleParser() *hackerNewsArticleParser {
	return &hackerNewsArticleParser{client: getClient()}
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func takeRandomArticesIds(ids []int, take int) []int {
	if take == 0 {
		return make([]int, 0)
	}
	if take >= len(ids) {
		return ids
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomIds := make([]int, 0, take)
	for i := 0; i < take; i++ {
		index := r.Intn(len(ids))
		randomIds = append(randomIds, ids[index])
	}

	return randomIds
}

func getTopArticlesIds(client *http.Client) ([]int, error) {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "dev-news-bot")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while parsing hackernews top articless, status: %s", resp.Status)
	}

	defer resp.Body.Close()
	var res []int
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	l := len(res)
	log.Infof("Parsed %d ids from hackernews", l)
	return res, nil
}

func getArticle(id int, client *http.Client) (*hackernewsArticle, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "dev-news-bot")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while parsing hackernews article by id: %d, status: %s", id, resp.Status)
	}

	defer resp.Body.Close()
	var res hackernewsArticle
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *hackerNewsArticleParser) Parse(ctx context.Context) ([]model.Article, error) {
	result := make([]model.Article, 0)
	ids, err := getTopArticlesIds(p.client)
	if err != nil {
		return nil, err
	}
	ids = takeRandomArticesIds(ids, 10)
	for _, id := range ids {
		hackerNewsArticle, err := getArticle(id, p.client)
		if err != nil {
			log.WithField("id", id).WithError(err).Error("error while parsing article by id")
			continue
		}
		article := model.NewArticle(hackerNewsArticle.Title, hackerNewsArticle.URL)
		if article.IsValid() {
			result = append(result, model.NewArticle(hackerNewsArticle.Title, hackerNewsArticle.URL))
		}
	}
	return result, nil
}
