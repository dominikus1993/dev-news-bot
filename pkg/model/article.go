package model

import (
	"net/url"
	"time"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/go-toolkit/crypto"
	"github.com/dominikus1993/go-toolkit/random"
	"github.com/samber/lo"
)

type ArticleId = string

type Article struct {
	id        ArticleId
	title     string
	content   string
	link      string
	source    string
	crawledAt time.Time
}

func (article Article) GetID() ArticleId {
	return article.id
}

func (article Article) GetTitle() string {
	return article.title
}

func (article Article) GetContent() string {
	return article.content
}

func (article Article) GetSource() string {
	return article.source
}

func (article Article) GetLink() string {
	return article.link
}

func (article Article) GetCrawledAt() time.Time {
	return article.crawledAt
}

type ArticlesStream <-chan Article

func NewArticleWithContent(title, link, content, source string) Article {
	id, _ := crypto.GenerateId(title, link)
	return Article{
		id:        id,
		title:     title,
		content:   content,
		link:      link,
		source:    source,
		crawledAt: time.Now().UTC(),
	}
}

func NewArticle(title, link, source string) Article {
	id, _ := crypto.GenerateId(title, link)
	return Article{
		id:        id,
		title:     title,
		content:   "",
		link:      link,
		source:    source,
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
	return random.TakeRandomToSlice(stream, take)
}

func UniqueArticles(articles ArticlesStream) ArticlesStream {
	return channels.UniqBy(articles, func(el Article) ArticleId { return el.id }, 10)
}

func UniqueArticlesArray(articles []Article) []Article {
	return lo.UniqBy(articles, func(article Article) ArticleId {
		return article.id
	})
}
