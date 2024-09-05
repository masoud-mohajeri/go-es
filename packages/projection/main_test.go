package projection

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"

	"github.com/masoud-mohajeri/go-es/packages/queue"
)

type MockQueue struct {
	frames []queue.Frame
}

func (m *MockQueue) Insert(frame queue.Frame) {
	// TODO implement if needed
	panic("implement me")
}

func (m *MockQueue) Show() {
	// TODO implement if needed
	panic("implement me")
}

func (m *MockQueue) Pop() (queue.Frame, bool) {
	if len(m.frames) == 0 {
		return queue.Frame{}, false
	}
	event := m.frames[0]
	m.frames = m.frames[1:]
	return event, true
}

func TestNewProjection(t *testing.T) {
	handlers := Handler{
		"event1": func(payload string) {
			fmt.Println("Handler for event1 called with payload:", payload)
		},
	}
	mockQueue := &MockQueue{}
	queues := []queue.FifoQueue{mockQueue}

	p := NewProjection(handlers, queues)

	if !assert.Equal(t, p.handlerMap, handlers) {
		t.Errorf("Expected handler map to be %v, got %v", handlers, p.handlerMap)
	}
	if len(p.queues) != 1 || p.queues[0] != mockQueue {
		t.Errorf("Expected queues to contain the mock queue, got %v", p.queues)
	}
}

func TestGenerateSnapshot(t *testing.T) {
	mockQueue := &MockQueue{
		frames: []queue.Frame{
			{Decoder: "event1", Payload: "test_payload"},
		},
	}

	// Use a WaitGroup to track the handler invocation
	var wg sync.WaitGroup
	wg.Add(1)

	handlers := Handler{
		"event1": func(payload string) {
			// Assert: Check the payload is correct
			if payload != "test_payload" {
				t.Errorf("Expected payload 'test_payload', got %s", payload)
			}
			wg.Done() // Mark the handler as called
		},
	}

	p := NewProjection(handlers, []queue.FifoQueue{mockQueue})
	p.GenerateSnapshot()

	wg.Wait()

	if len(mockQueue.frames) != 0 {
		t.Errorf("Expected queue to be empty after snapshot, but it has %d events left", len(mockQueue.frames))
	}
}

func TestGenerateSnapshot_EmptyQueue(t *testing.T) {
	mockQueue := &MockQueue{}
	handlers := Handler{
		"event1": func(payload string) {
			t.Errorf("Handler should not be called for an empty queue")
		},
	}

	p := NewProjection(handlers, []queue.FifoQueue{mockQueue})
	p.GenerateSnapshot()

	if len(mockQueue.frames) != 0 {
		t.Errorf("Expected queue to remain empty, but it has %d events", len(mockQueue.frames))
	}
}
