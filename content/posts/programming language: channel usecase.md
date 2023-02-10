---
title: "Channel Use Case in Golang"
date: 2022-12-27
authors: ["ngoctd"]
draft: false
tags: ["basic", "golang"]
---

This article will show many channel use cases

- Asynchronous and concurrency programing with Go channels is easy and enjoyable.
- The channel synchronization technique has wider range of uses and have more variables than the synchronization solutions used in some other languages, such as the actor model and the async/await pattern.

## Use Channels as Futures/Promises

**Return receive-only channels as results**

```go
func longTimeRequest() <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		// simulate workload run 3s in using goroutine
		time.Sleep(time.Second * 3)
		ch <- struct{}{}
	}()
    // return immediately

	return ch
}

func main() {
	now := time.Now()

	ch1 := longTimeRequest()
	ch2 := longTimeRequest()

	_, _ = <-ch1, <-ch2 // get result in future

	log.Println("exec in: ", time.Since(now)) // ~ 3s
}
```

**Pass send-only channels as arguments**

```go
func longTimeRequest(ch chan<- struct{}) {
	// simulate workload run 3s in using goroutine

	ch <- struct{}{}
	log.Println("send value to channel")
}

func main() {
	now := time.Now()
	// buffer = 2 to avoid block to handle channle receive
	ch := make(chan struct{}, 2)

	go longTimeRequest(ch)
	go longTimeRequest(ch)

	<-ch
	<-ch

	log.Println("exec in: ", time.Since(now))
}
```

**The first response win**

Sometimes, a piece of data can be received from several sources to avoid high latencies. For a lot of factors, the response durations of these sources may vary much. To make the response duration as short as possible, we can send a request to every source in separated goroutine. Only the first response use case will be used, other slower ones will be discarded.

```go
func source(c chan<- struct{}) {
	c <- struct{}{}
}

func main() {
	// c must be a buffered channel
	// If there are N sources, the capacity of the communication channel must be at least N-1
	// to avoid the goroutines corresponding the discarded responses being blocked for ever
	c := make(chan struct{}, 5)
	for i := 0; i < cap(c); i++ {
		go source(c)
	}

	// first response win
	<-c
}
```

## Use Channels for Notifications

Notifications can be viewed as special requests/responses in which the responded values are not important. Generally, we use the blank type struct{} as the element types of the notification channels.

**1-to-1 notification by sending a value to a channel**

If there are no values to be received from a channel, then the next receive operation on the channel will block until another goroutine sends a value to the channel. So we can send a value to a channel to notify another goroutine which is waiting to receive a value from the same channel.

```go
func blocking(c chan<- struct{}) {
	time.Sleep(1 * time.Second)
	// unlock by notifycation
	c <- struct{}{}
}

func main() {
	now := time.Now()
	ch := make(chan struct{})
	go blocking(ch)

	// blocking
	<-ch
	log.Println("since: ", time.Since(now))

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

```

**1-to-1 notification by receiving a value from a channel**

If the value buffer queue of a channel is full (the buffer queue of an unbuffered channel is always full), a send operation on the channel will block until another goroutine receives a value from the channel. So we can receive a value from a channel to notify another goroutine which is waiting to send a value to the same channel.

```go
func blocking(c <-chan struct{}) {
	time.Sleep(1 * time.Second)
	// unblock the second send in main goroutine
	<-c
}

func main() {
	now := time.Now()
	ch := make(chan struct{})
	go blocking(ch)

	// blocked here, wait for a notification
	ch <- struct{}{}
	log.Println("since: ", time.Since(now))

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
```

**N-to-1 and 1-to-N notifications**

By making using of the feature that infinite values can be received from a closed channel, we can close a channel to broadcast notifications (1-to-N 1 main to N worker). Use sync.WaitGroup for N-to-1(N worker to 1 main) notification

```go
func worker(ready chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done() // N-to-1 notification
	<-ready         // block here and wait notification
	time.Sleep(time.Second)
	log.Println("run")
}

func main() {
	ready := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go worker(ready, wg)
	go worker(ready, wg)
	go worker(ready, wg)

	close(ready)
	wg.Wait()
}
```

**Timer: Scheduled notification**

A custom one-time timer notification

```go
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
```

## Use Channels as Mutex Locks

There are two manners to use one-capacity buffered channels as mutex locks

1. Lock through a send, unlock through a receive

```go
func main() {
	cnt := 0
	wg := sync.WaitGroup{}

	lock := make(chan struct{}, 1)
	var increase = func() {
		defer wg.Done()
		// lock through send
		lock <- struct{}{}
		cnt++
		// unlock
		<-lock
	}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go increase()
	}

	wg.Wait()
	log.Println(cnt)

}

```

2. Lock through a receive, unlock through a send

```go
func main() {
	cnt := 0
	wg := sync.WaitGroup{}

	lock := make(chan struct{}, 1)
	lock <- struct{}{}

	var increase = func() {
		defer wg.Done()
		// lock through receive
		<-lock
		cnt++
		// unlock
		lock <- struct{}{}
	}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go increase()
	}

	wg.Wait()
	log.Println(cnt)

}
```

## Use Channels as Counting Semaphores

Buffered channels can be used as counting semaphores
Counting semaphores can be viewed as multi-owner locks. If the capacity of a channel is N, then it can be viewed as a lock which can have most N owners at any time. 

Counting semaphores are often used to enforce a maximum number of concurrent requests.

There are two manners to acquire one piece of ownership of a channel semaphore
1. Acquire ownership through a send, release through a receive
2. Acquire ownership through a receive, release through a send

