package usecase

import (
	"context"
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
)

type ArticleDto struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Link      string    `json:"link"`
	CrawledAt time.Time `json:"crawledAt"`
}

func NewArticleDto(article model.Article) ArticleDto {
	return ArticleDto{
		Title:     article.Title,
		Content:   article.Content,
		Link:      article.Link,
		CrawledAt: article.CrawledAt,
	}
}

type GetArticlesUseCase struct {
	articlesReader repositories.ArticlesReader
}

func NewGetArticlesUseCase(articlesReader repositories.ArticlesReader) *GetArticlesUseCase {
	return &GetArticlesUseCase{articlesReader: articlesReader}
}

func (u *GetArticlesUseCase) Execute(ctx context.Context, cmd repositories.GetArticlesParams) ([]ArticleDto, error) {
	result, err := u.articlesReader.Read(ctx, cmd)

	if err != nil {
		return make([]ArticleDto, 0), err
	}

	articles := make([]ArticleDto, len(result))
	for i, article := range result {
		articles[i] = NewArticleDto(article)
	}
	return articles, nil
}
