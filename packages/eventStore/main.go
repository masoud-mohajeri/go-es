package eventStore

import "sync"

type Storage interface {
	Get(key string) ([]byte, error)
	Save(key string, value []byte) error
}

func NewEventStore() Storage {
	return &txtFileStorage{
		mutex: sync.Mutex{},
	}
}
