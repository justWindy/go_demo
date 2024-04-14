package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/marusama/cyclicbarrier"
)

type FizzBuzz struct {
	n       int
	barrier cyclicbarrier.CyclicBarrier
	wg      sync.WaitGroup
}

func New(n int) *FizzBuzz {
	return &FizzBuzz{
		n:       n,
		barrier: cyclicbarrier.New(4),
	}
}

func (fb *FizzBuzz) start() {
	fb.wg.Add(4)
	go fb.fizz()
	go fb.buzz()
	go fb.fizzbuzz()
	go fb.number()

	fb.wg.Wait()
}

func (fb *FizzBuzz) fizz() {
	defer fb.wg.Done()
	ctx := context.Background()
	v := 0
	for {
		fb.barrier.Await(ctx)
		v++
		if v > fb.n {
			return
		}
		if v % 3 == 0 {
			if v % 5 == 0 {
				continue
			}
			if v == fb.n {
				fmt.Print(" fizz. ")
			} else {
				fmt.Print(" fizz,")
			}
		}
	}
}

func (fb *FizzBuzz) buzz() {
	defer fb.wg.Done()
	ctx := context.Background()
	v := 0
	for {
		fb.barrier.Await(ctx)
		v++
		if v > fb.n {
			return
		}
		if v%5 == 0 {
			if v%3 == 0 {
				continue
			}
			if v == fb.n {
				fmt.Print(" buzz. ")
			} else {
				fmt.Print(" buzz,")
			}
		}
	}
}

func (fb *FizzBuzz) fizzbuzz() {
	defer fb.wg.Done()
	ctx := context.Background()
	v := 0 
	for {
		fb.barrier.Await(ctx)
		v++
		if v > fb.n {
			return
		}
		if v%5 == 0 && v%3 == 0 {
			if v == fb.n {
				fmt.Print(" fizzbuzz. ")
			} else {
				fmt.Print(" fizzbuzz,")
			}
		}
	}
}

func (fb *FizzBuzz) number() {
	defer fb.wg.Done()
	ctx := context.Background()
	v := 0
	for {
		fb.barrier.Await(ctx)
		v++
		if v > fb.n {
			return
		}
		if v%5 != 0 && v%3 != 0 {
			if v == fb.n {
				fmt.Printf(" %d. ", v)
			} else {
				fmt.Printf(" %d,", v)
			}
		}
	}
}

func main() {
	fb := New(15)
	fb.start()
}
