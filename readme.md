# Basic Go Task Scheduler

[![Total alerts](https://img.shields.io/lgtm/alerts/g/jpollard-cs/GoTaskScheduler.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/jpollard-cs/GoTaskScheduler/alerts/)

## running the example

to run demo (runs for 10 seconds) just use Run button in repl.it or run the following in the console:

```bash
go run .
```

## running the tests

to run tests (add optional `-race` to check for race conditions):

```bash
go test -v -race
```

## what's next?

This is a super basic task scheduler written in Go

A thread safe min heap is utilized with a convenience method to get all jobs scheduled to run before a certain time

This is the first bit of code I've written in Go so it has some rough edges right now

It is also is a bit of an adjustment from other languages I've used in the past so there may be some things I haven't quite used in the right way, but I'm looking forward to learning how to better leverage the language

Some things that could be improved on (that I'm aware of):
- better error handling especially around channels
- improved testing coverage of safeminheap and the scheduler
- rather than using a ticker scheduling a timer that can be reset if a task with an earlier execution is scheduled (but with enough tasks perhaps a ticker could be more efficient?)
- use a wait group to keep the scheduler alive (`defer wg.Done()` before `for { select { ... } }` which will exit when cancelled), use a buffered channel as a semaphore to limit concurrency of how many tasks are run at once to help prevent resource exhaustion
- without much a much more complex design something like this has limited utility beyond code you're fully in control over (tasks could spawn more threads, use lots of memory, tasks are not sandboxed so there are security concerns)
