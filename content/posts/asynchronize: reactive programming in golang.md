---
title: "Reactive Programming with Go"
date: 2022-11-26
authors: ["ngoctd"]
image: "/images/reactive/rx.png"
draft: false
tags: ["reactive", "go"]
---

ReactiveX, or Rx for short, is an API for programming with Observable streams. ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises, and deferred. It is about processing streams of events of items, with events being any occurrences or changes within the system. A stream of events is called an Observable.

## RxGO

```go
```

## Hot vs. Cold Observables

In the RX world, there is a distinction between cold and hot Observables. When the data is produced by the Observable itself, it is a cold Observable. When the data is produced outside the Observable, it is a hot Observable. Usually, we don't want to create a producer over and over again, we favor a hot Observable.

```go
// Create a hot Observable using FromChannel operator
ch := make(chan rxgo.Item)
go func() {
    for i := 0; i < 3; i++ {
        ch <- rxgo.Of(i)
    }
    close(ch)
}()

observable := rxgo.FromChannel(ch)

// First Observer
for item := range observable.Observe() {
    fmt.Println("First Observer", item.V)
}

// Second Observer
for item := range observable.Observe() {
    fmt.Println("Second Observer", item.V)
}


// Result
//
// 0
// 1
// 2
// It means the first Observer already consumed all items.
```

```go
// Create a cold Observable using Defer operator:
observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
    for i := 0; i < 3; i++ {
        ch <- rxgo.Of(i)
    }
}})

// In the case of a cold observable, the stream was created independent for
// every Observer
for item := range observable.Observe() {
    fmt.Println("first observable: ", item.V)
}

// In the case of a cold observable, the stream was created independent for
// every Observer
for item := range observable.Observe() {
    fmt.Println("second observable: ", item.V)
}

// Result
// 0 
// 1
// 2
// 0
// 1
// 2

```

Hot vs cold Observable are not about how you consume items, it's about where data is produced. Good example for hot Observable are price ticks from a trading exchange. And if you teach an Observable to fetch products from a database, then yield them one by one, you will create the cold Observable.