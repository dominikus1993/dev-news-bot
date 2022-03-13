package console

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/k0kubun/pp"
)

type PPrintNotifier struct {
}

func NewPPrintNotifier() *PPrintNotifier {
	return &PPrintNotifier{}
}

func (not *PPrintNotifier) Notify(ctx context.Context, articles []model.Article) error {
	scheme := pp.ColorScheme{
		Integer: pp.Green | pp.Bold,
		Float:   pp.Black | pp.BackgroundWhite | pp.Bold,
		String:  pp.Yellow,
	}

	// Register it for usage
	pp.SetColorScheme(scheme)
	pp.Println("Parsed: ")
	for _, article := range articles {
		pp.Println(article)
	}
	return nil
}
