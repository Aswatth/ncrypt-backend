package utils

import "sync"

type Event struct {
	Type string
	Data map[string]string
}

type EventHandler func(data map[string]string) error

type Broadcast struct {
	mutex_lock sync.RWMutex
	subcribers map[string][]EventHandler
}

func NewBroadcast() *Broadcast {
	return &Broadcast{subcribers: make(map[string][]EventHandler)}
}

func (b *Broadcast) Subscribe(event_type string, handler EventHandler) {
	b.mutex_lock.Lock()
	defer b.mutex_lock.Unlock()
	b.subcribers[event_type] = append(b.subcribers[event_type], handler)
}

func (b *Broadcast) Publish(event Event) {
	b.mutex_lock.RLock()
	defer b.mutex_lock.RUnlock()

	if handlers, exists := b.subcribers[event.Type]; exists {
		var wg sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			go func(h EventHandler, data map[string]string) {
				defer wg.Done()
				h(data)
			}(handler, event.Data)
		}
		wg.Wait()
	}
}
