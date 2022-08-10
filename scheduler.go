package main

import (
	"errors"
	"fmt"
	"time"
)

type Scheduler struct {
	PQ         *SafeMinHeap
	scheduleCh chan Job
	workCh     chan Job
	timer      *time.Ticker
	cancelCh   chan struct{}
}

// convenience method for making the scheduler with some
// reasonable defaults
func MakeAndInitScheduler() Scheduler {
	s := Scheduler{}
	s.Init(
		MakeNewSafeMinHeap(0),
		make(chan Job),
		make(chan Job),
		make(chan struct{}),
		50,
	)
	return s
}

func (scheduler *Scheduler) Init(
	priorityQ *SafeMinHeap,
	scheduleCh chan Job,
	workCh chan Job,
	cancelCh chan struct{},
	tickDelayInMs int64,
) {
	scheduler.PQ = priorityQ
	scheduler.scheduleCh = scheduleCh
	scheduler.workCh = workCh
	scheduler.cancelCh = cancelCh
	scheduler.timer = time.NewTicker(time.Millisecond * time.Duration(tickDelayInMs))
	go func() {
		for {
			select {
			// TODO: it may make sense to have the listener for
			// the scheduler channel run in a separate goroutine
			// to ensure time-sensitive tasks are processed in a
			// timely manner
			case newJob, open := <-scheduler.scheduleCh:
				if !open {
					return
				}
				fmt.Println("scheduling job")
				go func() {
					if err := scheduler.PQ.Add(newJob); err != nil {
						// right now we'd only get here if PQ has a size
						// limit that we hit
						// a blocking queue doesn't make much
						// sense in the context of scheduling
						// in reality a persistent data store
						// or for ephemeral time sensitive work
						// a more advanced in-memory system
						// (Redis supports priority queues)
						// would be needed to improve on this
						panic(err)
					}
				}()
			case jobToRun, open := <-scheduler.workCh:
				if !open {
					return
				}
				go func() {
					fmt.Printf(
						"running job %s at %v with scheduled time %v\r\n",
						jobToRun.name,
						time.Now(),
						jobToRun.scheduledTime,
					)
					jobToRun.task()
					// TODO report completion, use "semaphore" via buffered channel to limit concurrency
				}()
			// this will drop ticks if we for some reason slow down
			// which is useful as it won't create a backlog
			case t := <-scheduler.timer.C:
				fmt.Println("timer ticked")
				go func() {
					fmt.Println("checking for jobs to run")
					jobsToRun := scheduler.PQ.PopBefore(t)
					fmt.Printf("found %d jobs to run\r\n", len(jobsToRun))
					if jobsToRun == nil || len(jobsToRun) == 0 {
						return
					}
					for i := 0; i < len(jobsToRun); i++ {
						scheduler.workCh <- jobsToRun[i]
					}
				}()
			}
		}
	}()

	return
}

func (scheduler *Scheduler) Schedule(task func(), delayInMs int) {
	job := Job{
		task:          task,
		scheduledTime: time.Now().Add(time.Duration(delayInMs) * time.Millisecond),
	}
	scheduler.scheduleCh <- job
}

// Shuts down the scheduler
// if Scheduler is already canceled then returns canceled error
func (scheduler *Scheduler) Cancel() error {
	select {
	case <-scheduler.cancelCh:
		return errors.New("scheduler already canceled")
	default:
		fmt.Println("canceling scheduler")
		scheduler.timer.Stop()
		if len(scheduler.timer.C) > 0 {
			// drain channel
			<-scheduler.timer.C
		}
    // TODO: still a potential
    // race condition here
    // if timer ticked prior to stopping
    // and resulting goroutine
    // (e.g. maybe was blocked on priority queue)
    // calls to do work after
    // work channel is closed
		close(scheduler.cancelCh)
		close(scheduler.workCh)
		close(scheduler.scheduleCh)
	}
	return nil
}
