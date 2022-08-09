package main

import (
	"testing"
	"time"
)

func TestMinHeap(t *testing.T) {
	times := make([]time.Time, 10)
	times[0] = time.Now()
	for i := 1; i < 10; i++ {
		times[i] = times[i-1].Add(time.Second)
	}

	t.Run("no limit min heap", func(t *testing.T) {
		h := MakeNewSafeMinHeap(0)
		h.Add(Job{scheduledTime: times[5]})
		h.Add(Job{scheduledTime: times[4]})
		h.Add(Job{scheduledTime: times[6]})
		h.Add(Job{scheduledTime: times[3]})
		h.Add(Job{scheduledTime: times[4]})

		result := []int{3, 4, 4, 5, 6}
		for _, i := range result {
			if peek := h.Peek(); !peek.scheduledTime.Equal(times[i]) {
				t.Errorf("peek = %v expected %v", peek.scheduledTime, times[i])
			}
			if pop := h.Pop(); !pop.scheduledTime.Equal(times[i]) {
				t.Errorf("pop = %v expected %v", pop.scheduledTime, times[i])
			}
		}
		if h.Size() != 0 {
			t.Errorf("expect empty but size = %v", h.Size())
		}
		if peek := h.Peek(); !peek.scheduledTime.IsZero() {
			t.Errorf("empty peek = %v", peek.scheduledTime)
		}
		if pop := h.Pop(); !pop.scheduledTime.IsZero() {
			t.Errorf("empty pop = %v", pop.scheduledTime)
		}
	})

	t.Run("limited min heap", func(t *testing.T) {
		max := 5
		h := MakeNewSafeMinHeap(max)
		h.Add(Job{scheduledTime: times[5]})
		h.Add(Job{scheduledTime: times[4]})
		h.Add(Job{scheduledTime: times[6]})
		h.Add(Job{scheduledTime: times[3]})
		h.Add(Job{scheduledTime: times[4]})

		if err := h.Add(Job{scheduledTime: times[1]}); err != ErrMax {
			t.Errorf("max add expected error: %v, but %v", ErrMax, err)
		}

		result := []int{3, 4, 4, 5, 6}
		for _, i := range result {
			if peek := h.Peek(); !peek.scheduledTime.Equal(times[i]) {
				t.Errorf("peek = %v expected %v", peek.scheduledTime, times[i])
			}
			if pop := h.Pop(); !pop.scheduledTime.Equal(times[i]) {
				t.Errorf("pop = %v expected %v", pop.scheduledTime, times[i])
			}
		}
		if h.Size() != 0 {
			t.Errorf("expect empty but size = %v", h.Size())
		}
		if peek := h.Peek(); !peek.scheduledTime.IsZero() {
			t.Errorf("empty peek = %v", peek.scheduledTime)
		}
		if pop := h.Pop(); !pop.scheduledTime.IsZero() {
			t.Errorf("empty pop = %v", pop.scheduledTime)
		}
	})

	t.Run("pop many", func(t *testing.T) {
		max := 5
		h := MakeNewSafeMinHeap(max)
		h.Add(Job{scheduledTime: times[4]})
		h.Add(Job{scheduledTime: times[3]})
		h.Add(Job{scheduledTime: times[2]})
		h.Add(Job{scheduledTime: times[1]})
		h.Add(Job{scheduledTime: times[0]})

		tasks := h.PopBefore(times[3])

		if len(tasks) != 3 {
			t.Errorf("expected 3 tasks")
		}

		for i, task := range tasks {
			if !task.scheduledTime.Equal(times[i]) {
				t.Errorf("got task %d with scheduled time %v, expected scheduled time to be %v", i, task.scheduledTime, times[i])
			}
		}
	})
}
