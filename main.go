package main

func main() {
}

func isClose(c chan int) bool {
	select {
	case _, ok := <-c:
		return ok
	default:
	}

	return false
}
