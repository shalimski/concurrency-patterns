package main

import (
	"fmt"
	"sync"
)

type PubSub[T any] struct {
	subscribers []chan T
	mu          sync.RWMutex
	closed      bool
	closeOnce   sync.Once
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: nil,
		mu:          sync.RWMutex{},
		closed:      false,
	}
}

func (ps *PubSub[T]) Subscribe() <-chan T {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil
	}

	subscriber := make(chan T)
	ps.subscribers = append(ps.subscribers, subscriber)
	return subscriber
}

func (ps *PubSub[T]) Publish(v T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, subscriber := range ps.subscribers {
		subscriber <- v
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.closeOnce.Do(func() {
		ps.closed = true

		for _, subscriber := range ps.subscribers {
			close(subscriber)
		}
	})

}

func (ps *PubSub[T]) Closed() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.closed
}

func main() {

	ps := NewPubSub[int]()

	wg := sync.WaitGroup{}

	subscriber1 := ps.Subscribe()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case v, ok := <-subscriber1:
				if !ok {
					fmt.Println("subscriber1 was closed")
					return
				}

				fmt.Printf("new value for subscriber1: %d\n", v)

			}
		}
	}()

	subscriber2 := ps.Subscribe()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case v, ok := <-subscriber2:
				if !ok {
					fmt.Println("subscriber2 was closed")
					return
				}

				fmt.Printf("new value for subscriber2: %d\n", v)

			}
		}
	}()

	ps.Publish(1)
	ps.Publish(100)
	ps.Close()

	wg.Wait()

}
