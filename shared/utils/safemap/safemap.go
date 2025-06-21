package safemap

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

func (s *SafeMap[K, V]) Put(key K, value V) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	val, ok := s.m[key]
	s.mu.Unlock()
	return val, ok
}

func (s *SafeMap[K, V]) Contains(key K) bool {
	_, ok := s.Get(key)
	return ok
}

func (s *SafeMap[K, V]) Delete(key K) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}
