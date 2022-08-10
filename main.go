package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	s := MakeAndInitScheduler()
	go func() {
		for {
			s.Schedule(func() {}, rand.Intn(1000))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		time.Sleep(time.Second * 10)
		s.Cancel()
		cancel()
	}()
	<-ctx.Done()
	fmt.Println(ctx.Err() == context.Canceled)
}
