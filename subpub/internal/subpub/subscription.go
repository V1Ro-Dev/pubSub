package subpub

import (
	"log"
	"sync"

	"github.com/google/uuid"

	"pubsub/utils"
)

type SubscriptionImpl struct {
	id       uuid.UUID
	subject  string
	handler  MessageHandler
	messages chan interface{}
	isActive bool
	mu       sync.RWMutex
}

func NewSubscription(subject string, handler MessageHandler) *SubscriptionImpl {
	return &SubscriptionImpl{
		id:       utils.GenerateID(),
		subject:  subject,
		handler:  handler,
		messages: make(chan interface{}, 100),
		isActive: true,
		mu:       sync.RWMutex{},
	}
}

func (s *SubscriptionImpl) Unsubscribe() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isActive {
		s.isActive = false
		log.Printf("%s sub unsubscribed from %s", s.id, s.subject)
		close(s.messages)
	}
}

func (s *SubscriptionImpl) Listen(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for msg := range s.messages {
		s.mu.RLock()
		if !s.isActive {
			s.mu.RUnlock()
			return
		}
		s.mu.RUnlock()

		log.Printf("%s sub got notification from %s: %v", s.id, s.subject, msg)
		s.handler(msg)
	}
}

func (s *SubscriptionImpl) getMessage(msg interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isActive {
		select {
		case s.messages <- msg:

		default:
			log.Println("sub channel full")
		}
	}
}
