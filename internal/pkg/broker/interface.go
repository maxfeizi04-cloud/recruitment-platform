package broker

import (
	"context"
	"fmt"
	"sync"
)

type MessageBroker interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(topic string, handler func(ctx context.Context, payload []byte) error) error
	Close() error
}

type InMemoryBroker struct {
	handlers map[string][]func(context.Context, []byte) error
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewInMemoryBroker() *InMemoryBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &InMemoryBroker{
		handlers: make(map[string][]func(context.Context, []byte) error),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (b *InMemoryBroker) Publish(ctx context.Context, topic string, payload []byte) error {
	b.mu.RLock()
	handlers, ok := b.handlers[topic]
	b.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no subscribers for topic: %s", topic)
	}

	for _, handler := range handlers {
		go func(h func(context.Context, []byte) error) {
			if err := h(ctx, payload); err != nil {
				fmt.Printf("[broker] handler error for topic %s: %v\n", topic, err)
			}
		}(handler)
	}

	return nil
}

func (b *InMemoryBroker) Subscribe(topic string, handler func(ctx context.Context, payload []byte) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}

func (b *InMemoryBroker) Close() error {
	b.cancel()
	return nil
}
