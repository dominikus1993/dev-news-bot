package model

import (
	"net/url"
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
	_, err := url.Parse(a.Link)
	if err != nil {
		return false
	}
	contentIsValid := len(a.Content) == 0 || len(a.Content) < 2048
	titleIsValid := a.Title != ""
	return contentIsValid && titleIsValid && linkIsValid
}
