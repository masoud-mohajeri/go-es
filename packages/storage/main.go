package storage

import "sync"

type Storage interface {
	Get(key string) ([]byte, error)
	Save(key string, value []byte) error
}

func NewStorage() Storage {
	return &txtFileStorage{
		mutex: sync.Mutex{},
	}
}
