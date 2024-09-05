package dispatcher

import "github.com/masoud-mohajeri/go-es/packages/queue"

type EventDispatcher struct {
	queue queue.FifoQueue // I need it to be a reference, so I can pass it to different dispatchers the same q
}

type Dispatcher interface {
	Dispatch(payload string, kind string) //  better typing for kind ???
}

func (c *EventDispatcher) Dispatch(payload string, kind string) {
	c.queue.Insert(queue.Frame{
		Payload: payload,
		Decoder: kind,
	})
}

func NewDispatcher(queue queue.FifoQueue) Dispatcher {
	return &EventDispatcher{
		queue: queue,
	}
}
