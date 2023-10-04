package utils

import (
	"math"
	"sync"
)

type Queue struct {
	sync.Mutex
	maxItems int
	Items    []any
}

func NewQueue() *Queue {
	q := &Queue{}
	q.maxItems = math.MaxInt
	return q
}

func NewQueueLimited(maxCapacity int) *Queue {
	q := &Queue{}
	q.maxItems = maxCapacity
	return q
}

func (q *Queue) Count() int {
	q.Lock()
	defer q.Unlock()
	return len(q.Items)
}

func (q *Queue) Space() int {
	q.Lock()
	defer q.Unlock()
	return q.maxItems - len(q.Items)
}

// Push if there is space
func (q *Queue) MaybePush(item any) bool {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) >= q.maxItems {
		return false
	}
	q.Items = append(q.Items, item)
	return true
}

// Push no matter what
func (q *Queue) Push(item any) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, item)
}

func (q *Queue) Pop() any {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}
