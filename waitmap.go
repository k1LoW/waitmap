package waitmap

import "sync"

type WaitMap[K comparable, V any] struct {
	lockmap map[K]*sync.Cond
	valmap  map[K]V
	mu      sync.Mutex
}

// New creates a new WaitMap.
func New[K comparable, V any]() *WaitMap[K, V] {
	return &WaitMap[K, V]{
		lockmap: map[K]*sync.Cond{},
		valmap:  map[K]V{},
	}
}

// Get returns the value associated with the key.
// If the key does not exist, Get blocks until the key is set.
func (m *WaitMap[K, V]) Get(key K) V {
	m.mu.Lock()
	lock, ok := m.lockmap[key]
	if !ok {
		lock = sync.NewCond(&m.mu)
		m.lockmap[key] = lock
	}
	m.mu.Unlock()
	for {
		m.mu.Lock()
		if v, ok := m.valmap[key]; ok {
			m.mu.Unlock()
			return v
		}
		lock.Wait()
		m.mu.Unlock()
	}
}

// TryGet returns the value associated with the key.
// If the key does not exist, TryGet returns (zero value, false).
func (m *WaitMap[K, V]) TryGet(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.valmap[key]
	return v, ok
}

// Set sets the value associated with the key.
func (m *WaitMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[key]
	if !ok {
		lock = sync.NewCond(&m.mu)
		m.lockmap[key] = lock
	}
	m.valmap[key] = value
	lock.Broadcast()
}

// Delete deletes the value associated with the key.
func (m *WaitMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[key]
	if ok {
		lock.Broadcast()
		delete(m.lockmap, key)
	}
	delete(m.valmap, key)
}
