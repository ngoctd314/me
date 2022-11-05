package main

import (
	"fmt"
	"time"
)

func generator(data string) <-chan string {
	channel := make(chan string)

	go func() {
		for {
			channel <- data
			time.Sleep(time.Duration(100 * time.Millisecond))
		}
	}()
	return channel
}

func main() {
	ch := generator("ngoctd")
	go func() {
		for v := range ch {
			fmt.Println(v)
		}
	}()
	time.Sleep(time.Second * 5)
}
