package notifications

import (
	"context"
	"sync"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	log "github.com/sirupsen/logrus"
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

func NewBroadcaster(notifiers []Notifier) *broadcaseter {
	return &broadcaseter{notifiers: notifiers}
}

func (b *broadcaseter) Broadcast(ctx context.Context, articles []model.Article) error {
	var wg sync.WaitGroup
	for _, notifier := range b.notifiers {
		wg.Add(1)
		go func(ctx context.Context, n Notifier, wait *sync.WaitGroup) {
			defer wait.Done()
			err := n.Notify(ctx, articles)
			if err != nil {
				log.WithError(err).Error("Error while broadcasting")
			}
		}(ctx, notifier, &wg)
	}
	wg.Wait()
	return nil
}
