package model

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
	contentIsValid := len(a.Content) > 0 && len(a.Content) < 2048
	titleIsValid := a.Title != ""
	linkIsValid := a.Link != ""
	// uri, err := url.Parse(a.Link)
	// if err != nil {
	// 	return false
	// }
	return contentIsValid && titleIsValid && linkIsValid
}
