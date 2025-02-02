package events

import (
	"sync"
)

var instance *EventManager

// When emiting an event this is the struct that will come through
type Event struct {
	Name string
	Data any
}

// EventHandler defines the function that will handle an emitted event
type EventHandler func(event Event)

// For now this is a "dumb" event manager in that there is no way to keep track of one handler from another
// it may be beneficial in the future to change this to a system that allows for unsubscribes
type EventManager struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// EventManager is a singleton so there should be no way to instantiate one directly
func GetEventManagerInstance() *EventManager {
	if instance == nil {
		instance = &EventManager{
			subscribers: make(map[string][]EventHandler),
		}
	}

	return instance
}

// This will add a handler to the slice of handlers for an event name
func (em *EventManager) Subscribe(eventName string, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.subscribers[eventName] = append(em.subscribers[eventName], handler)
}

// This allows for abitrary event dispatch and each handler for that event will be notified
func (em *EventManager) Emit(eventName string, data any) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	event := Event{
		Name: eventName,
		Data: data,
	}

	if handlers, exists := em.subscribers[eventName]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}
