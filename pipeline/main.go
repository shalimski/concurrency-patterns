package main

import "fmt"

func main() {
	for v := range square(square(generator(1, 2, 3, 4))) {
		fmt.Printf("%d\n", v)
	}
}

func square(c <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range c {
			out <- v * v
		}
		close(out)
	}()

	return out
}

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v
		}
	}()

	return out
}
