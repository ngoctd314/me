package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan int, 1)
	for i := 0; i < 5; i++ {
		go source(ctx, c)
	}
	rnd := <-c // only the first response win
	cancel()
	log.Println("winner: ", rnd)
	time.Sleep(time.Second * 5)
	fmt.Println("number of goroutines:", runtime.NumGoroutine())
}

func source(ctx context.Context, c chan<- int) {
	now := time.Now()
	fn := func(ctx context.Context) <-chan int {
		defer func() {
			log.Println("makesure goroutines is released")
		}()
		rand.Seed(time.Now().UnixNano())

		fmt.Println("call to source")
		// try-send
		ch := make(chan int, 1)

		ra, rb := rand.Int(), rand.Intn(3)+1
		select {
		// simulate work load, we don't need to call when context is canceled
		case <-time.After(time.Duration(rb) * time.Second):
			defer func() {
				log.Println("end call to source after:", time.Since(now))
			}()
			ch <- ra
		case <-ctx.Done():
		}

		return ch
	}

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	case v := <-fn(ctx):
		select {
		case c <- v:
		default:
		}
	}

}
