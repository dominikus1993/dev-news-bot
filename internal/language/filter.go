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
	languages := []lingua.Language{
		lingua.English,
		lingua.Polish,
	}
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	return languageFilter{detector: detector}
}

func (filter languageFilter) Where(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		for article := range articles {
			if lang, exists := filter.detector.DetectLanguageOf(article.GetTitle()); exists {
				fmt.Println(lang)
				result <- article
			}
		}
		close(result)
	}()
	return result
}
