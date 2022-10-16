package language

import (
	"context"
	"fmt"

	"github.com/dominikus1993/dev-news-bot/pkg/filters"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/pemistahl/lingua-go"
)

type languageFilter struct {
	detector lingua.LanguageDetector
}

func NewLanguageFilter() filters.ArticlesFilter {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()

	return languageFilter{detector: detector}
}

func getArticleTitleAndContent(article model.Article) string {
	content := article.GetContent()
	if content == "" {
		return article.GetTitle()
	}
	return fmt.Sprintf("%s %s", article.GetTitle(), content)
}

func (filter languageFilter) Where(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		for article := range articles {
			if lang, exists := filter.detector.DetectLanguageOf(getArticleTitleAndContent(article)); exists {
				if lang == lingua.English || lang == lingua.Polish {
					result <- article
				}
			}
		}
		close(result)
	}()
	return result
}
