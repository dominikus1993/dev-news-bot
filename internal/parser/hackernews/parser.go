package hackernews

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/go-toolkit/random"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
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
	client              *http.Client
	maxArticlesQuantity int
}

func NewHackerNewsArticleParser(maxArticlesQuantity int) *hackerNewsArticleParser {
	return &hackerNewsArticleParser{client: getClient(), maxArticlesQuantity: maxArticlesQuantity}
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func takeRandomArticesIds(ids []int, take int) []int {
	return random.TakeRandomFromSlice(ids, take)
}

func getTopArticlesIds(client *http.Client) ([]int, error) {
	const url = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
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

func (p *hackerNewsArticleParser) Parse(ctx context.Context) model.ArticlesStream {
	result := make(chan model.Article, 20)
	go func() {
		defer close(result)
		ids, err := getTopArticlesIds(p.client)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorln("Error while parsing hackernews top articles")
			return
		}
		ids = takeRandomArticesIds(ids, p.maxArticlesQuantity)
		for _, id := range ids {
			hackerNewsArticle, err := getArticle(id, p.client)
			if err != nil {
				log.WithField("id", id).WithError(err).Errorln("error while parsing article by id")
				continue
			}
			article := model.NewArticle(hackerNewsArticle.Title, hackerNewsArticle.URL, source)
			if article.IsValid() {
				result <- article
			}
		}
	}()
	return result
}
