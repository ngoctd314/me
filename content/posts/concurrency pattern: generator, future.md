---
title: "Generator and Future Pattern"
date: 2022-11-05
authors: ["ngoctd"]
draft: false
tags: ["concurrency pattern", "golang"]
series: ["concurrency"]
---

Generator Pattern allows the consumer of the data produced by the generator to run in parallel when the generator function is busy computing the next value. 

A Future indicates any data that is needed in future but its computation can be started in parallel so that it can be fetched from the background when needed.

## Generator

Generator Pattern is used to generator a sequence of values which is used to produce some output. This allows the consumer of the data produced by the generator to run in parallel when the generator function is busy computing the text value.

```go
func fib(n int) chan int {
	c := make(chan int)

	go func() {
        // next state is depend on previous state
		for i, j := 0, 1; i < n; i, j = i+j, i {
			c <- i
		}
		close(c)
	}()

	return c
}

func consumer(c chan int) {
    for v := range c {
        fmt.Println(v)
    }
}
```

The generator and the consumer can work concurrently (maybe in parallel) as the logic involved in both are different.

## Future

A Future indicates any data that is needed in future but its computation can be started in parallel so that it can be fetched from the background when needed. Mostly, futures are used to send asynchronous http request.

```go
type data struct {
    Body []byte
    Error error
}

func futureDate(url string) <- chan data {
    c := make(chan data, 1)

    go func(){
        resp, err := http.Get(url)
        if err != nil {
            c <- data{
            	Body:  nil,
            	Error: err,
            }
            return
        }

        body, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()

        if err != nil {
            c <- data{
            	Body:  nil,
            	Error: err,
            }
            return
        }

        c <- data{
        	Body:  body,
        	Error: err,
        }
    }()

    return nil
}
```

The actual http request is done asynchronously in a goroutine. The main function can continue doing other things. When the result is needed, we read the result from the channel. If the request isn't finished yet, the channel will block until the result is ready.

## Different Between Generator and Future 

In generator pattern, we generate next state base on previous state (maybe not), but it purpose is compute many things in background. In future pattern, we use goroutine to execute an "heavy job" (only one job).