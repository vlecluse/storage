package storage 

import (
	"sync"
)

type Queue struct {
	sync.RWMutex
	capacity uint32
	queue []uint32
	set map[uint32]bool
}

func newQueue(capacity uint32) *Queue {
	return &Queue{
		capacity: capacity,
		queue: make([]uint32, 0, capacity),
		set: make(map[uint32]bool),
	}
}

func (q *Queue) enqueue(index uint32) {
	if index < 0 || index >= q.capacity {
		return
	}

	q.RLock()
	isset := q.set[index]
	q.RUnlock()

	if isset {
		return
	}

	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, index)
	q.set[index] = true
}

func (q *Queue) dequeue() (uint32, bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return 0, false
	}

	index := q.queue[0]
	q.queue = q.queue[1:]
	q.set[index] = false

	return index, true
}

func (q *Queue) shiftToBack() (uint32, bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return 0, false
	}

	index := q.queue[0]
	q.queue = append(q.queue[1:], index)
	return index, true
}

func (q *Queue) len() int {
	q.RLock()
	defer q.RUnlock()

	return len(q.queue)
}

func (q *Queue) remove(index uint32) { 
	q.Lock()
	defer q.Unlock()
    newLen := 0

    // Iterate through the slice and keep the elements that are not equal to the target value
    for _, v := range q.queue {
        if v != index {
            q.queue[newLen] = v
            newLen++
        }
    }

	q.queue = q.queue[:newLen]
	q.set[index] = false
}