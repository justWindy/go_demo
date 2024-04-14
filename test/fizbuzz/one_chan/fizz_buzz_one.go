package main

import (
	"fmt"
	"sync"
)

type FizzBuzz struct {
	n  int
	ch chan int
	wg sync.WaitGroup
}

func New(n int) *FizzBuzz {
	return &FizzBuzz{
		n:  n,
		ch: make(chan int, 1),
	}
}

func (fb *FizzBuzz) fizz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%3 == 0 {
			if v%5 == 0 {
				fb.ch <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" fizz. ")
			} else {
				fmt.Print(" fizz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) buzz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 == 0 {
			if v%3 == 0 {
				fb.ch <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" buzz. ")
			} else {
				fmt.Print(" buzz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) fizzbuzz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 == 0 && v%3 == 0 {
			if v == fb.n {
				fmt.Print(" fizzbuzz.")
			} else {
				fmt.Print(" fizzbuzz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) number() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 != 0 && v%3 != 0 {
			if v == fb.n {
				fmt.Printf(" %d. ", v)
			} else {
				fmt.Printf(" %d,", v)
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) start() {
	fb.wg.Add(4)
	go fb.fizz()
	go fb.buzz()
	go fb.fizzbuzz()
	go fb.number()

	fb.ch <- 1
	fb.wg.Wait()
}

func main() {
	fb := New(15)
	fb.start()
}
