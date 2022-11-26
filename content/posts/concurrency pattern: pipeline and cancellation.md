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

**References**
- https://go.dev/blog/pipelines