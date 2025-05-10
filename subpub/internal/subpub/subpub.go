package subpub

import (
	"context"
	"sync"
)

type SubPubImpl struct {
	storage map[string][]*SubscriptionImpl
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

func NewPubSub() SubPubImpl {
	return SubPubImpl{
		storage: make(map[string][]*SubscriptionImpl),
		mu:      sync.RWMutex{},
		wg:      sync.WaitGroup{},
	}
}

func (sp *SubPubImpl) Publish(subject string, msg interface{}) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if subs, ok := sp.storage[subject]; ok {
		for _, sub := range subs {
			sub.getMessage(msg)
		}
	}

	return nil
}

func (sp *SubPubImpl) Subscribe(subject string, handler MessageHandler) (Subscription, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sub := NewSubscription(subject, handler)
	sp.storage[subject] = append(sp.storage[subject], sub)

	go sub.Listen(&sp.wg)

	return sub, nil
}

func (sp *SubPubImpl) Close(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		sp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
