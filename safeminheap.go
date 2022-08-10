package main

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

// a Job designed for the MinHeap
// TODO: consider adding a priority and updating
// the Less method to account for this
type Job struct {
	name          string    // the name of the job
	task          func()    // the task to execute for the job
	scheduledTime time.Time // the scheduled time to run the job
	index         int       // The index of the job in the heap.
}

type JobHeap []Job

// Len is length of JobHeap
func (h *JobHeap) Len() int {
	return len(*h)
}

// Less means Job j is newer than i
func (h *JobHeap) Less(i, j int) bool {
	return !(*h)[i].scheduledTime.After((*h)[j].scheduledTime)
}

// Swap swaps the elements with indexes i and j.
func (h *JobHeap) Swap(i, j int) {
	if h.Len() > 0 {
		(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	}
}

// Push adds x to tail
func (h *JobHeap) Push(x interface{}) {
	item, ok := x.(Job)
	if !ok {
		return
	}
	*h = append((*h)[:], item)
}

// Pop removes x from head
func (h *JobHeap) Pop() (x interface{}) {
	if h.Len() > 0 {
		l := h.Len() - 1
		x = (*h)[l]
		(*h)[l] = Job{}
		*h = (*h)[:l]
	}
	return
}

// errors
var (
	ErrMax = errors.New("heap max size")
)

// convenience wrapper around heap
type SafeMinHeap struct {
	mu   sync.Mutex
	heap JobHeap
	max  int
}

func MakeNewSafeMinHeap(max int) *SafeMinHeap {
	if max < 0 {
		max = 0
	}
	h := &SafeMinHeap{heap: make(JobHeap, 0, max), max: max}
	heap.Init(&h.heap)
	return h
}

func (h *SafeMinHeap) Add(j Job) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.max > 0 && len(h.heap) >= h.max {
		return ErrMax
	}
	heap.Push(&h.heap, j)
	return nil
}

// consider using PopAfter instead
func (h *SafeMinHeap) Pop() Job {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.heap) == 0 {
		return Job{}
	}
	return heap.Pop(&h.heap).(Job)
}

// while this peek is thread safe
// it's marked unsafe since
// it can't guarantee you peek what you pop
// consider using PopAfter
func (h *SafeMinHeap) Peek() Job {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.heap) == 0 {
		return Job{}
	}
	return h.heap[0]
}

// pops Jobs scheduled before the provided time
func (h *SafeMinHeap) PopBefore(t time.Time) []Job {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.heap) == 0 {
		return nil
	}

	jobs := make([]Job, 0, 5)
	for len(h.heap) > 0 && h.heap[0].scheduledTime.Before(t) {
		j := heap.Pop(&h.heap).(Job)
		jobs = append(jobs, j)
	}
	return jobs
}

func (h *SafeMinHeap) Size() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.heap)
}
