package usecase

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/providers"
)

type parseArticlesAndSendItUseCase struct {
	articlesProvider providers.ArticlesProvider
}

func NewParseArticlesAndSendItUseCase() *parseArticlesAndSendItUseCase {
	return &parseArticlesAndSendItUseCase{}
}

func (u *parseArticlesAndSendItUseCase) Execute(ctx context.Context, articlesQuantity int) error {
	_, err := u.articlesProvider.Provide(ctx)
	if err != nil {
		return err
	}
	return nil
}
