package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	now := time.Now()

	queue := 10
	ok := make(chan int, queue)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < queue; i++ {
		go reply(ctx, ok, i)
	}

	// ai nhanh thì thắng
	firstRepsonse := <-ok
	cancel()
	log.Println("since: ", time.Since(now), "pos", firstRepsonse)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func reply(ctx context.Context, ch chan<- int, pos int) {
	rd := rand.Intn(3) + 1
	select {
	case <-time.After(time.Duration(rd) * time.Second):
		log.Println("run")
	case <-ctx.Done():
		log.Println("cancel")
	}

	ch <- pos
}
