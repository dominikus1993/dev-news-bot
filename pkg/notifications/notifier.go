package notifications

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"golang.org/x/sync/errgroup"
)

type Notifier interface {
	Notify(ctx context.Context, articles []model.Article) error
}

type Broadcaster interface {
	Broadcast(ctx context.Context, articles []model.Article) error
}

type broadcaseter struct {
	notifiers []Notifier
}

func NewBroadcaster(notifiers ...Notifier) *broadcaseter {
	return &broadcaseter{notifiers: notifiers}
}

func (b *broadcaseter) Broadcast(ctx context.Context, articles []model.Article) error {
	wg, ctx := errgroup.WithContext(ctx)
	for _, notifier := range b.notifiers {
		not := notifier
		wg.Go(func() error {
			return not.Notify(ctx, articles)
		})
	}
	return wg.Wait()
}
