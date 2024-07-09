package waitmap

import "sync"

type waitMap[K comparable, V any] struct {
	lockmap map[K]*sync.Cond
	valmap  map[K]V
	mu      sync.Mutex
}

// New creates a new waitMap.
func New[K comparable, V any]() *waitMap[K, V] {
	return &waitMap[K, V]{
		lockmap: make(map[K]*sync.Cond),
		valmap:  make(map[K]V),
	}
}

// Get returns the value associated with the key k.
// If the key does not exist, Get blocks until the key is set.
func (m *waitMap[K, V]) Get(k K) V {
	m.mu.Lock()
	lock, ok := m.lockmap[k]
	if !ok {
		lock = sync.NewCond(&m.mu)
		m.lockmap[k] = lock
	}
	m.mu.Unlock()
	for {
		m.mu.Lock()
		if v, ok := m.valmap[k]; ok {
			m.mu.Unlock()
			return v
		}
		lock.Wait()
		m.mu.Unlock()
	}
}

// TryGet returns the value associated with the key k.
// If the key does not exist, TryGet returns (zero value, false).
func (m *waitMap[K, V]) TryGet(k K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.valmap[k]
	return v, ok
}

// Set sets the value associated with the key k.
func (m *waitMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[k]
	if !ok {
		lock = sync.NewCond(&m.mu)
		m.lockmap[k] = lock
	}
	m.valmap[k] = v
	lock.Broadcast()
}

// Delete deletes the value associated with the key k.
func (m *waitMap[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[k]
	if ok {
		lock.Broadcast()
		delete(m.lockmap, k)
	}
	delete(m.valmap, k)
}
