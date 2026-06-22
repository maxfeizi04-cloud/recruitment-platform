package broker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/broker"
)

func TestInMemoryBroker_PublishSubscribe(t *testing.T) {
	b := broker.NewInMemoryBroker()
	defer b.Close()

	var mu sync.Mutex
	var received []string

	err := b.Subscribe("test.topic", func(ctx context.Context, payload []byte) error {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, string(payload))
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	err = b.Publish(context.Background(), "test.topic", []byte("hello"))
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Errorf("received count = %d, want 1", len(received))
	}
	if received[0] != "hello" {
		t.Errorf("received[0] = %s, want hello", received[0])
	}
}

func TestInMemoryBroker_NoSubscriber(t *testing.T) {
	b := broker.NewInMemoryBroker()
	defer b.Close()

	err := b.Publish(context.Background(), "no.such.topic", []byte("data"))
	if err == nil {
		t.Error("expected error for topic with no subscribers")
	}
}
