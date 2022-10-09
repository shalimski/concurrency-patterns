package main

import "fmt"

func main() {
	c := generator(1, 2, 3, 4, 5)

	for i := range c {
		fmt.Printf("%d\n", i)
	}
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
