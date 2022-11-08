package main

func main() {

}

func generatePipeline(numbers []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, number := range numbers {
			out <- number
		}
		close(out)
	}()

	return out
}

func squareNumber(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}

func fanIn()
