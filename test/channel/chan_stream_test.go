package channel

import (
	"fmt"
	"testing"

	"golang.org/x/exp/constraints"
)

func asStream[T any](done <-chan struct{}, values ...T) <-chan T {
	s := make(chan T)
	go func() {
		defer close(s)

		for _, v := range values {
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()

	return s
}

func takeN[T any](done <-chan struct{}, valueStream <-chan T, num int) <-chan T {
	takeStream := make(chan T)
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func sqrt[T constraints.Integer](in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func mapChan[T, K any](in <-chan T, fn func(T) K) <-chan K {
	out := make(chan K)
	if in == nil {
		close(out)
		return out
	}

	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()

	return out
}

func reduce[T, K any](in <-chan T, fn func(r K, v T) K) K {
	var out K

	if in == nil {
		return out
	}

	for v := range in {
		out = fn(out, v)
	}

	return out
}

func TestMapReduce(t *testing.T) {
	in := make(chan int, 2)
	go func() {
		defer close(in)
		in <- 10
		in <- 5
	}()

	mapFn := func(v int) int {
		return v * 10
	}

	reduceFn := func(r, v int) int {
		return r + v
	}

	sum := reduce(mapChan(in, mapFn), reduceFn)
	fmt.Println(sum)
}
