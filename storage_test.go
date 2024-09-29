package storage

import(
	"testing"

	"github.com/stretchr/testify/assert"
)

const storageSize = 128 * 1024
const key = uint64(1)
const modulus = 1 * storageSize

var value = []byte("Test")
var mediumValue = createValueWithSize(64*1024)
var bigValue = createValueWithSize(256*1024)

func createValueWithSize(size int) []byte {
	var buf []byte
	for i := 0; i < size; i++ {
		buf = append(buf, byte(i))
	}
	return buf
}

func getNewStorage() *Storage[[]byte] {
	return NewMemoryStorage[[]byte](storageSize, 8)
}

func TestSetAndGet(t *testing.T) {
	t.Parallel()

	storage := getNewStorage()

	storage.Set(uint64(1), value)
	result, _ := storage.Get(uint64(1))
	assert.Equal(t, value, result)
	_, exists := storage.Get(uint64(2))
	assert.Equal(t, false, exists)

	storageSizeToUse := storageSize
	if storageSizeToUse < MinimumSize {
		storageSizeToUse = MinimumSize
	}

	valueSize := len(bigValue)
	maxSize := valueSize * storageSizeToUse

	for i := 0; i < 80 * 1024; i++ {
		storage.Set(uint64(i), bigValue)
	}

	storedSize := storage.getSize()

	assert.GreaterOrEqual(t, maxSize, storedSize)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	storage := getNewStorage()
	storage.Set(uint64(1), value)
	result, _ := storage.Get(uint64(1))
	assert.Equal(t, value, result)

	storage.Delete(uint64(2))
	result, _ = storage.Get(uint64(1))
	assert.Equal(t, value, result)

	storage.Delete(uint64(1))
	_, exists := storage.Get(uint64(1))
	assert.Equal(t, false, exists)
}

func BenchmarkStorage_Set(b *testing.B) {
	storage := getNewStorage()
	b.RunParallel(func(pb *testing.PB) {
		keyHash := uint64(1)
		for pb.Next() {
			storage.Set(keyHash%modulus, value)
			keyHash++
		}
	})
}

func BenchmarkStorage_SetMedium(b *testing.B) {
	storage := getNewStorage()

	b.RunParallel(func(pb *testing.PB) {
		keyHash := uint64(1)
		for pb.Next() {
			storage.Set(keyHash%modulus, mediumValue)
			keyHash++
		}
	})
}

func BenchmarkStorage_SetBig(b *testing.B) {
	storage := getNewStorage()

	b.RunParallel(func(pb *testing.PB) {
		keyHash := uint64(1)
		for pb.Next() {
			storage.Set(keyHash%modulus, bigValue)
			keyHash++
		}
	})
}

func BenchmarkStorage_Get(b *testing.B) {
	storage := getNewStorage()

	b.RunParallel(func(pb *testing.PB) {
		keyHash := uint64(1)
		for pb.Next() {
			storage.Get(keyHash%modulus)
			keyHash++
		}
	})
}

func BenchmarkStorage_Delete(b *testing.B) {
	storage := getNewStorage()

	b.RunParallel(func(pb *testing.PB) {
		keyHash := uint64(1)
		for pb.Next() {
			storage.Delete(keyHash%modulus)
			keyHash++
		}
	})
}
