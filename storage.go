package storage

const MinimumSize = 10 * 1024

type Storage[T any] struct {
	shard uint64
	memoryShard []*MemoryShard[T]
}

func NewMemoryStorage[T any](size int, shardPower int) *Storage[T] {
	shard := 1 << shardPower

	storage := &Storage[T]{
		memoryShard: make([]*MemoryShard[T], shard),
		shard: uint64(shard),
	}

	if size < MinimumSize {
		size = MinimumSize
	}

	segmentSize := size / shard

	for i := 0; i < shard; i++ {
		storage.memoryShard[i] = newMemoryShard[T](segmentSize)
	}

	return storage
}

func (s *Storage[T]) Set(key uint64, value T) {
	segmentId := s.getSegmentId(key)
	s.memoryShard[segmentId].set(key, value)
}

func (s *Storage[T]) Get(key uint64) (T, bool) {
	segmentId := s.getSegmentId(key)
	return s.memoryShard[segmentId].get(key)
}

func (s *Storage[T]) Delete(key uint64) {
	segmentId := s.getSegmentId(key)
	s.memoryShard[segmentId].delete(key)
}

func (m *Storage[T]) getSegmentId(key uint64) int {
	return int(key % m.shard)
}

func (s *Storage[T]) getSize() int {
	var size int

	for segmentId := range s.memoryShard {
		size += s.memoryShard[segmentId].getSize()
	}

	return size
}