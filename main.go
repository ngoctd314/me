package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	c1 := crawl("crawl1", time.Second)
	c2 := crawl("crawl2", time.Second*5)
	c := fanIn(c1, c2)

	// forward to another source
	for i := 0; i < 10; i++ {
		log.Println(<-c)
	}

}
func crawl(id string, sleepTime time.Duration) chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("crawl for %s: %d", id, i)
			time.Sleep(sleepTime)
		}
	}()

	return c
}

func fanIn(c ...chan string) <-chan string {
	fi := make(chan string)

	for _, v := range c {
		go func(v chan string) {
			for i := range v {
				fi <- i
			}
		}(v)
	}

	return fi
}
