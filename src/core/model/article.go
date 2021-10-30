package model

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

type Article struct {
	Title   string
	Content string
	Link    string
}

func NewArticleWithContent(title, link, content string) Article {
	return Article{
		Title:   title,
		Content: content,
		Link:    link,
	}
}

func NewArticle(title, link string) Article {
	return Article{
		Title:   title,
		Content: "",
		Link:    link,
	}
}

func (a *Article) IsValid() bool {
	linkIsValid := a.Link != ""
	if !linkIsValid {
		return false
	}
	link, err := url.ParseRequestURI(a.Link)
	if err != nil {
		return false
	}
	fmt.Println(link.Scheme)
	contentIsValid := len(a.Content) == 0 || len(a.Content) < 2048
	titleIsValid := a.Title != ""
	return contentIsValid && titleIsValid
}

func GetRandomArticles(articles []Article, take int) []Article {
	if take == 0 {
		return make([]Article, 0)
	}
	if take >= len(articles) {
		return articles
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomArticles := make([]Article, 0, take)
	for _, i := range r.Perm(take) {
		randomArticles = append(randomArticles, articles[i])
	}

	return randomArticles
}
