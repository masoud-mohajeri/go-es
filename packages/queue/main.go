package queue

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Frame struct {
	Id      string
	Payload string // json value of dispatcher object
	Decoder string
}

// FifoQueue this is a queue with short polling pattern
// TODO: create a queue with long polling pattern -> rabbitMQ like behavior
// TODO: create a queue with observer pattern -> kafka like behavior
type FifoQueue interface {
	Pop() (Frame, bool)
	Insert(frame Frame)
	Show()
}

type Queue struct {
	mutex  *sync.Mutex
	frames []Frame
}

func (q *Queue) Pop() (frame Frame, ok bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.frames) == 0 {
		return Frame{}, false
	}

	poppedFrame := q.frames[0]
	q.frames = q.frames[1:]

	return poppedFrame, true

}

func (q *Queue) Insert(frame Frame) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	frame.Id = uuid.NewString()
	q.frames = append(q.frames, frame)
}

func (q *Queue) Show() {
	frames := "Frames: "
	for _, frame := range q.frames {
		frames += frame.Decoder + ", "
	}
	fmt.Println(frames)
}

func NewQueue() FifoQueue {
	return &Queue{
		mutex:  new(sync.Mutex),
		frames: make([]Frame, 0),
	}
}
