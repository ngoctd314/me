package main

import (
	"log"
	"sync"
)

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
