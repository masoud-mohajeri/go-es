package projection

import (
	"fmt"
	"github.com/masoud-mohajeri/go-es/packages/queue"
)

type Handler map[string]func(payload string)

type Projection struct {
	// TODO: add (mutex property + Register method) for both
	handlerMap Handler // map
	queues     []queue.FifoQueue
}

func (p *Projection) GenerateSnapshot() {
	fmt.Println("-projection: Generating snapshot from", len(p.queues), "queues")
	for _, q := range p.queues {
		go func() {
			event, ok := q.Pop()
			if !ok {
				fmt.Println("-projection: Queue is empty")
				return
			}
			if handler, ok := p.handlerMap[event.Decoder]; ok {
				handler(event.Payload)
			}
			// we are using go 1.22, so we don't need to pass q as arg to this func invocation :D
		}()
	}

}

func NewProjection(handlers Handler, queues []queue.FifoQueue) Projection {
	return Projection{
		handlerMap: handlers,
		queues:     queues,
	}
}
