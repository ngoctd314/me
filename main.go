package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// The first stage, gen is a function that converts a list of integers to a channel
// that emits the integers in the list. The gen function starts a goroutine that sends
// the integers on the channel and closes the channel when all the values have been sent
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		// close when all the values have been sent
		close(out)
	}()

	return out
}

// The second stage, sq receives integers from a channel and returns a channel that emits the square
// of each received integer. After the inbound channel is closed and this stage has sent all the values
// downstream, it closes the outbound channel
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// consume inbound
		for v := range in {
			out <- v * v
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, v := range cs {
		go output(v)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func worker(ready chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done() // N-to-1 notification
	<-ready         // block here and wait notification
	time.Sleep(time.Second)
	log.Println("run")
}

func after(duration time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(duration)
		ch <- struct{}{}
	}()

	return ch
}

func main() {
	now := time.Now()
	<-after(time.Second)
	fmt.Println(time.Since(now))
}
