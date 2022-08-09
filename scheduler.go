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

func MakeAndInitScheduler(pq *SafeMinHeap) Scheduler {
	s := Scheduler{
		PQ: pq,
	}
	s.Init(make(chan Job), make(chan Job), make(chan struct{}), 50)
	return s
}

func (scheduler *Scheduler) Init(
	scheduleCh chan Job,
	workCh chan Job,
	cancelCh chan struct{},
	tickDelayInMs int64,
) {
	scheduler.scheduleCh = scheduleCh
	scheduler.workCh = workCh
	scheduler.cancelCh = cancelCh
	scheduler.timer = time.NewTicker(time.Millisecond * time.Duration(tickDelayInMs))
	go func() {
		for {
			select {
			case <-scheduler.cancelCh:
				return
			case newJob := <-scheduler.scheduleCh:
				fmt.Println("scheduling job")
				go func() {
					if err := scheduler.PQ.Add(newJob); err != nil {
						panic(err)
					}
				}()
			case jobToRun := <-scheduler.workCh:
				go func() {
					fmt.Printf(
						"running job %s at %v with scheduled time %v\r\n",
						jobToRun.name,
						time.Now(),
						jobToRun.scheduledTime,
					)
					jobToRun.task()
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

// Close shutdown scheduler and workers goroutine.
// if Scheduler is already closed then returns ErrClosed.
func (scheduler *Scheduler) Cancel() error {
	select {
	case <-scheduler.cancelCh:
		return errors.New("scheduler already canceled")
	default:
		fmt.Println("canceling scheduler")
		close(scheduler.cancelCh)
		close(scheduler.workCh)
		close(scheduler.scheduleCh)
	}
	return nil
}
