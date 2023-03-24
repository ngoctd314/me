package main

import (
	"fmt"
	"sync"
)

func main() {
	n := 10
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go notify(wg)
	}
	wg.Wait()
}

func notify(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("send notify")
}
