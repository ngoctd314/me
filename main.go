package main

import (
	"fmt"
	"sync"
)

func main() {
	for i := 0; i < 1e9; i++ {
		if fn() == "x: 0y: 0" {
			fmt.Println("RUNNx0,y0")
		}
		if fn() == "y: 0x: 0" {
			fmt.Println("RUNNy0,x0")
		}
	}

}

func fn() string {
	var (
		x, y int
		wg   sync.WaitGroup
		rs   = ""
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		x = 1
		rs += fmt.Sprintf("y: %d", y)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		y = 1

		rs += fmt.Sprintf("x: %d", x)
	}()

	wg.Wait()

	return rs
}

func fanIn()
