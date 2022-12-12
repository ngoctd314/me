---
title: "Pipelines and cancellation"
date: 2022-11-26
authors: ["ngoctd"]
draft: false
tags: ["concurrency pattern", "golang"]
series: ["concurrency"]
---

Go's concurrency primitives make it easy to construct streaming data pipelines that make efficient use of I/O and multiple CPUs. This article introduces what is pipeline, how to to construct a pipeline and introduces techniques for dealing with failures cleanly.

## What is a pipeline?

A pipeline is a series of stages connected by channels, where each channel is a group of goroutines running the same function. In each stage, the goroutines

- receive values from upstream via inbound channels 
- perform some function on that data, usually producing new values
- send the values downstream via outbound channels

Each stage has any number of inbound and outbound channels, except the first and last stages, which have only outbound or inbound channels. The first stage is sometimes called the source and producer; the last stage, the sink or consumer.

## Squaring numbers

Consider a pipeline with three stages

The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list.

```go
// the main function sets up the pipeline and runs the final stages
// it receives values from the second stage and prints each one, until the channel is closed
func main() {
	// a pipeline
	c := gen(2, 3)
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

// The gen function starts a goroutine that sends the integers on the channel and closes the channel when all the values have been sent:
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

// The second stage, sq receives integers from a channel and returns a channel that emits the square of each received integer.
// After the inbound channel is closed and this stage has sent all the values downstream, it closes the outbound channel
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
```
## Fan-in, Fan-out

## Stopping Short

There is  a pattern to our pipeline functions: 
- stages close their outbound channels when all the send operations are done.
- stages keep receiving value from inbound channels until those channel are closed.

This pattern allows each receiving stage to be written as for range loop and ensures that all goroutines exit once all values have been successfully sent downstream.

But in real pipelines, stages don't always receive all the inbound values. Sometimes this is by design: the receiver may only need a subset of values to make progress. More often, a stage exit early because an inbound value represents an error in a earlier stages to stop producing values that later stage don't need. 

```go
func gen(nums ...int) (result<-chan int) {
	result = make(chan int)
	for _, num := range nums {
		// blocking operation
		result <- num
	}

	return
}

func main() {
	result := gen(2,3)
	fmt.Println(result)

	// This is a resource leak: goroutines consume memory and runtime 
	// resources, and heap references in goroutine stacks heap data from
	// being garbage collected. Goroutines are not garbage collected;
	// they must exit on their own.
	fmt.Println(runtime.NumGoroutine()) // 2
}
```
We need to arrange for the upstream stages of our pipeline to exit even when the downstream stages fail to receive all the inbound values. One way to do this is to change the outbound channels to have a buffer.

```go
c := make(chan int, 2) // buffer size 2
c <- 1 // succeeds immediately
c <- 2 // succeeds immediately
c <- 3 // blocks until another goroutine does <- c and receives 1
```

When the **max** number of values to be sent is known at channel creation time, a buffer can simplify the code. 

## Explicit cancellation

When main decides to exit without receiving all the values from out, it must tell the goroutines in the upstream stages to abandon the values they're trying to send. It does so by sending values on channel called done.

```go
func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		for _, n := range nums {
			out <- n
		}
		// close when all the values have been sent
		close(out)
	}()

	return out
}

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

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	// we don't know number of channel size
	out := make(chan int)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(out)
	}()

	output := func(c <-chan int) {
		for n := range c {
			select {
			// consume data
			case out <- n:
			// stop send value
			case <-done:
				break
			}
		}
		wg.Done()
	}

	wg.Add(len(cs))
	// consume input channel, exit sending goroutine of input channel
	for _, v := range cs {
		// create goroutine to consume input channel data
		go output(v)
	}

	return out
}

func main() {
	// create main goroutine
	g := gen(2, 3) // create 2 new goroutines then exit

	c1 := sq(g) // create a new goroutine and block
	c2 := sq(g) // create a new goroutine and block
	// 3 goroutines because no goroutine consume it

	done := make(chan struct{})

	defer func() {
		close(done)
		// only main goroutine
	}()

	out := merge(done, c1, c2)
	fmt.Println(<-out)
}

func printGoroutine() {
	time.Sleep(time.Second)
	fmt.Println("num goroutine", runtime.NumGoroutine())
}
```

**Here are the guidelines for pipeline construction**
- Stages close their outbound channels when all the send operations are done.
- Stages keep receiving value from inbound channels until those channels are closed or the senders are unblocked.

## Digesting a tree



**References**
- https://go.dev/blog/pipelines