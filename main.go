package main

import (
	"fmt"
	"sync"
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

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in
	// fan-out
	c1 := sq(in)
	c2 := sq(in)

	// Consume the merged output from c1 and c2
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
