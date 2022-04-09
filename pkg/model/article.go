package model

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
)

type Article struct {
	id        string
	title     string
	content   string
	link      string
	crawledAt time.Time
}

func (article Article) GetID() string {
	return article.id
}

func (article Article) GetTitle() string {
	return article.title
}

func (article Article) GetContent() string {
	return article.content
}

func (article Article) GetLink() string {
	return article.link
}

func (article Article) GetCrawledAt() time.Time {
	return article.crawledAt
}

type ArticlesStream <-chan Article

func NewArticleWithContent(title, link, content string) Article {
	id, _ := GenerateId(title, link)
	return Article{
		id:        id,
		title:     title,
		content:   content,
		link:      link,
		crawledAt: time.Now().UTC(),
	}
}

func NewArticle(title, link string) Article {
	id, _ := GenerateId(title, link)
	return Article{
		id:        id,
		title:     title,
		content:   "",
		link:      link,
		crawledAt: time.Now().UTC(),
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (a *Article) IsValid() bool {
	linkIsValid := a.link != ""
	if !linkIsValid {
		return false
	}
	contentIsValid := len(a.content) == 0 || len(a.content) < 2048
	titleIsValid := a.title != ""
	return isUrl(a.link) && contentIsValid && titleIsValid
}

func TakeRandomArticles(stream ArticlesStream, take int) []Article {
	if take == 0 {
		return make([]Article, 0)
	}
	articles := channels.ToSlice(stream)
	if take >= len(articles) {
		return articles
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomArticles := make([]Article, 0, take)
	for i := 0; i < take; i++ {
		index := r.Intn(len(articles))
		randomArticles = append(randomArticles, articles[index])
	}

	return randomArticles
}

func UniqueArticles(articles ArticlesStream) ArticlesStream {
	seen := make(map[string]bool)
	res := make(chan Article, 20)
	go func() {
		for v := range articles {
			if !seen[v.id] {
				seen[v.id] = true
				res <- v
			}
		}
		close(res)
	}()
	return res
}
