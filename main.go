package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func blocking(c <-chan struct{}) {
	time.Sleep(1 * time.Second)
	// unblock the second send in main goroutine
	<-c
}

func main() {
	now := time.Now()
	ch := make(chan struct{}, 1)
	go blocking(ch)

	// blocked here, wait for a notification
	ch <- struct{}{}
	log.Println("since: ", time.Since(now))

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
