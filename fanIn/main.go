package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx := context.Background()

	a := generator(1, 2, 3, 4, 5)
	b := generator(6, 7, 8, 9, 10)
	c := generator(11, 12, 13, 14, 15)

	res := fanIn(ctx, a, b, c)
	for i := range res {
		fmt.Printf("%d\n", i)
	}
}

func fanIn(ctx context.Context, input ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(input))

	for _, ch := range input {
		ch := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}
					out <- v
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
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
