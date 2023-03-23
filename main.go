package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan string)
	go secret(ch)

	log.Println(<-ch)
}

func secret(ch chan<- string) {
	time.Sleep(time.Second)
	ch <- "this is secret"
}
