package waitmap

import "sync"

type waitMap[K comparable, V any] struct {
	lockmap map[K]chan struct{}
	valmap  map[K]V
	mu      sync.Mutex
}

func New[K comparable, V any]() *waitMap[K, V] {
	return &waitMap[K, V]{
		lockmap: make(map[K]chan struct{}),
		valmap:  make(map[K]V),
	}
}

func (m *waitMap[K, V]) Get(k K) V {
	m.mu.Lock()
	lock, ok := m.lockmap[k]
	if !ok {
		lock = make(chan struct{})
		m.lockmap[k] = lock
	}
	m.mu.Unlock()
	<-lock
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.valmap[k]
	if !ok {
		panic("key not found. this is a bug in the code.")
	}
	return v
}

func (m *waitMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[k]
	if !ok {
		lock = make(chan struct{})
		m.lockmap[k] = lock
	}
	var locked bool
	select {
	case _, locked = <-lock:
	default:
		locked = true
	}
	m.valmap[k] = v
	if locked {
		close(lock)
	}
}

func (m *waitMap[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.lockmap[k]
	if ok {
		delete(m.lockmap, k)
		select {
		case _, locked := <-lock:
			if locked {
				close(lock)
			}
		default:
		}
	}
	delete(m.valmap, k)
}
