package model

import (
	"math/rand"
	"net/url"
	"time"
)

type Article struct {
	ID        string
	Title     string
	Content   string
	Link      string
	CrawledAt time.Time
}

type ArticlesStream <-chan Article

func NewArticleWithContent(title, link, content string) Article {
	id, _ := GenerateId(title, link)
	return Article{
		ID:        id,
		Title:     title,
		Content:   content,
		Link:      link,
		CrawledAt: time.Now().UTC(),
	}
}

func NewArticle(title, link string) Article {
	id, _ := GenerateId(title, link)
	return Article{
		ID:        id,
		Title:     title,
		Content:   "",
		Link:      link,
		CrawledAt: time.Now().UTC(),
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (a *Article) IsValid() bool {
	linkIsValid := a.Link != ""
	if !linkIsValid {
		return false
	}
	contentIsValid := len(a.Content) == 0 || len(a.Content) < 2048
	titleIsValid := a.Title != ""
	return isUrl(a.Link) && contentIsValid && titleIsValid
}

func TakeRandomArticles(stream ArticlesStream, take int) []Article {
	if take == 0 {
		return make([]Article, 0)
	}
	articles := ToArticlesArray(stream)
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

func ToArticlesArray(s ArticlesStream) []Article {
	res := make([]Article, 0)
	for v := range s {
		res = append(res, v)
	}
	return res
}

func UniqueArticles(articles ArticlesStream) ArticlesStream {
	seen := make(map[string]bool)
	res := make(chan Article)
	go func() {
		for v := range articles {
			if !seen[v.ID] {
				seen[v.ID] = true
				res <- v
			}
		}
		close(res)
	}()
	return res
}
