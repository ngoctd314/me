package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

type Seat int
type Bar chan Seat

func (bar Bar) ServeCustomer(c int) {
	log.Print("customer#", c, " enters the bar")
	seat := <-bar // receive value from bar
	log.Print("++ customer#", c, "drinks at seat#", seat)
	time.Sleep(time.Second * time.Duration(10+rand.Intn(6)))
	log.Print("-- customer#", c, " free seat#", seat)
	bar <- seat

}

func main() {
	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar, 10)
	for seatId := 0; seatId < cap(bar24x7); seatId++ {
		bar24x7 <- Seat(seatId)
	}

	for customerId := 0; ; customerId++ {
		log.Println("num goroutine: ", runtime.NumGoroutine())
		time.Sleep(time.Second)
		go bar24x7.ServeCustomer(customerId)
	}

}
