---
title: "Fan out, Fan in"
date: 2022-11-05
authors: ["ngoctd"]
draft: true
tags: ["concurrency pattern", "golang"]
series: ["concurrency"]
---
Fan-in Fan-out is a way of Multiplexing and De-Multiplexing in golang. Fan-in refers to processing multiple input data and combining into a single entity. Fan-out is the exact opposite, dividing the data into multiple smaller chunks,   

You've got a pipeline set up. Data is flowing through your system beautifully. Sometimes, stages in your pipeline can be computationally expensive. When this happens, upstream stages in your pipeline can become blocked while waiting for your expensive stages to complete.

One of the interesting properties of pipelines is the ability they give you to operate on the stream of data using a combination of separate, often **reordered** stages. Maybe that would help improve the performance of the pipeline. In fact, it turns out it can, and this pattern has a name: fan-out, fan-in.

## Fan in,Fan out

Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline.

You might consider fanning out one of your stages if both of the following apply:
- It doesn't rely on values that the state had calculated before.
- It takes a long time to run. (system call, a heavy cpu job, ...)

The property of order-independence is important because you have no guarantee in what order concurrent copies of your stage will run, nor in what order they will return.


Multiple functions can be read from the same channel until that channel is closed; this is called fan-out. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.

A function can read from multiple inputs and proceed until all are closed by multiplexing then input channels onto a single channel that's closed when all the inputs are closed. This is called fan-in.

```go
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

```