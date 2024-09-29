package storage

import(
	"sync"
	"log"
)

type MemoryShard[T any] struct {
	sync.RWMutex
	items []T
	mapping sync.Map
	available *Queue
	fallback *Queue
}

func newMemoryShard[T any](segmentSize int) *MemoryShard[T]{
	memoryShard := &MemoryShard[T]{
		items: make([]T, segmentSize),
		available: newQueue(uint32(segmentSize)),
		fallback: newQueue(uint32(segmentSize)),
	}

	for i := uint32(0); i < uint32(segmentSize); i++ {
		memoryShard.available.enqueue(i)
	}

	return memoryShard
}

func (m *MemoryShard[T]) getIndex(key uint64) (uint32, bool) {
	index, exists := m.mapping.Load(key)

	if exists == false {
		return 0,false
	}

	return index.(uint32), exists
}

func (m *MemoryShard[T]) getAndSetIndex(key uint64) uint32 {
	index, exists := m.getIndex(key)

	if exists {
		return index
	}

	index = m.getAvailableIndex()
	m.mapping.Store(key, index)
	
	return index
}

func (m *MemoryShard[T]) updateMapping(key uint64, index uint32) {
	m.mapping.Store(key, index)
}

func (m *MemoryShard[T]) set(key uint64, value T) {
	index := m.getAndSetIndex(key)
	m.Lock()
	defer m.Unlock()
	m.items[index] = value
}

func (m *MemoryShard[T]) get(key uint64) (T, bool) {
	var emptyValue T
	index, exists := m.getIndex(key)

	if exists == false{
		return emptyValue, false
	}

	m.RLock()
	defer m.RUnlock()
	
	return m.items[index], true
}

func (m *MemoryShard[T]) delete(key uint64) {
	var emptyValue T
	index, exists := m.getIndex(key)
	if exists == false {
		return
	}

	m.Lock()
	m.items[index] = emptyValue
	m.Unlock()
	m.mapping.Delete(key)
	m.releaseIndex(index)
}

func (m *MemoryShard[T]) getSize() int {
	var size int

	for _, content := range m.items {
		value := any(content).([]byte)
		size = size + len(value)
	}

	return size
}

func (m *MemoryShard[T]) getAvailableIndex() uint32 {
	index, exists := m.available.dequeue()

	if exists {
		m.fallback.enqueue(index)
		return index
	}

	fallbackIndex, exists := m.fallback.shiftToBack()

	if exists {
		return fallbackIndex
	}

	log.Println("No index available")

	return 0
}

func (m *MemoryShard[T]) releaseIndex(index uint32) {
	m.available.enqueue(index)
	m.fallback.remove(index)
}