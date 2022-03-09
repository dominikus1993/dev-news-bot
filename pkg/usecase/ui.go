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

type ArticlesDto struct {
	Articles      []ArticleDto `json:"articles"`
	Total         int          `json:"total"`
	NumberOfPages int          `json:"numberOfPages"`
}

func countNumberOfPages(total int, pageSize int) int {
	if total%pageSize == 0 {
		return total / pageSize
	}
	return total/pageSize + 1
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

func (u *GetArticlesUseCase) Execute(ctx context.Context, cmd repositories.GetArticlesParams) (*ArticlesDto, error) {
	result, err := u.articlesReader.Read(ctx, cmd)

	if err != nil {
		return nil, err
	}

	articles := make([]ArticleDto, len(result.Articles))
	for i, article := range result.Articles {
		articles[i] = NewArticleDto(article)
	}
	return &ArticlesDto{Articles: articles, Total: result.Total, NumberOfPages: countNumberOfPages(result.Total, cmd.PageSize)}, nil
}
