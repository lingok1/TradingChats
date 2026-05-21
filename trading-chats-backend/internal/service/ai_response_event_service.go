package service

import (
	"sync"

	"trading-chats-backend/internal/models"
)

type AIResponseEventService struct {
	mutex       sync.RWMutex
	subscribers map[string]map[chan models.AIResponseEvent]struct{}
}

func NewAIResponseEventService() *AIResponseEventService {
	return &AIResponseEventService{
		subscribers: make(map[string]map[chan models.AIResponseEvent]struct{}),
	}
}

func (s *AIResponseEventService) Subscribe(tabTag string) (<-chan models.AIResponseEvent, func()) {
	normalizedTab := models.NormalizeTabTag(tabTag)
	ch := make(chan models.AIResponseEvent, 8)

	s.mutex.Lock()
	if _, ok := s.subscribers[normalizedTab]; !ok {
		s.subscribers[normalizedTab] = make(map[chan models.AIResponseEvent]struct{})
	}
	s.subscribers[normalizedTab][ch] = struct{}{}
	s.mutex.Unlock()

	unsubscribe := func() {
		s.mutex.Lock()
		if subscribers, ok := s.subscribers[normalizedTab]; ok {
			if _, exists := subscribers[ch]; exists {
				delete(subscribers, ch)
				close(ch)
			}
			if len(subscribers) == 0 {
				delete(s.subscribers, normalizedTab)
			}
		}
		s.mutex.Unlock()
	}

	return ch, unsubscribe
}

func (s *AIResponseEventService) Publish(event models.AIResponseEvent) {
	s.PublishToChannel(event.TabTag, event)
}

// PublishToChannel 发布到指定频道，event 内容不变（用于跨频道转发，例如 home 频道）
func (s *AIResponseEventService) PublishToChannel(channel string, event models.AIResponseEvent) {
	normalizedTab := models.NormalizeTabTag(channel)

	s.mutex.RLock()
	subscribers := s.subscribers[normalizedTab]
	for ch := range subscribers {
		select {
		case ch <- event:
		default:
		}
	}
	s.mutex.RUnlock()
}
