package lambda

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/usecase"
)

type ParseAndSendArticlesHandler struct {
	usecase *usecase.ParseArticlesAndSendItUseCase
}

func NewParseAndSendArticlesHandler(usecase *usecase.ParseArticlesAndSendItUseCase) *ParseAndSendArticlesHandler {
	return &ParseAndSendArticlesHandler{usecase: usecase}
}

func (h *ParseAndSendArticlesHandler) Handle(ctx context.Context) error {
	return h.usecase.Execute(ctx, 10)
}
