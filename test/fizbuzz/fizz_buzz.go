package main

import (
	"fmt"
	"sync"
)

type FizzBuzz struct {
	n   int
	chs []chan int
	wg  sync.WaitGroup
}

func NewFizzBuzz(n int) *FizzBuzz {
	chs := make([]chan int, 4)
	for i := 0; i < 4; i++ {
		chs[i] = make(chan int, 1)
	}
	return &FizzBuzz{n: n, chs: chs}
}

func (fb *FizzBuzz) start() {
	fb.wg.Add(4)
	go fb.fizz()
	go fb.buzz()
	go fb.fizzbuzz()
	go fb.number()
	fb.chs[0] <- 1
	fb.wg.Wait()
}

func (fb *FizzBuzz) fizz() {
	defer fb.wg.Done()
	next := fb.chs[1]
	for v := range fb.chs[0] {
		if v > fb.n {
			next <- v
			return
		}
		if v%3 == 0 {
			if v%5 == 0 {
				next <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" fizz. ")
			} else {
				fmt.Print(" fizz,")
			}
			next <- v + 1
			continue
		}
		next <- v
	}
}

func (fb *FizzBuzz) buzz() {
	defer fb.wg.Done()
	next := fb.chs[2]
	for v := range fb.chs[1] {
		if v > fb.n {
			next <- v
			return
		}
		if v%5 == 0 {
			if v%3 == 0 {
				next <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" buzz. ")
			} else {
				fmt.Print(" buzz,")
			}
			next <- v + 1
			continue
		}
		next <- v
	}
}

func (fb *FizzBuzz) fizzbuzz() {
	defer fb.wg.Done()
	next := fb.chs[3]
	for v := range fb.chs[2] {
		if v > fb.n {
			next <- v
			return
		}
		if v%5 == 0 && v%3 == 0 {
			if v == fb.n {
				fmt.Print(" fizzbuzz. ")
			} else {
				fmt.Print(" fizzbuzz,")
			}
			next <- v + 1
			continue
		}
		next <- v
	}
}

func (fb *FizzBuzz) number() {
	defer fb.wg.Done()
	next := fb.chs[0]
	for v := range fb.chs[3] {
		if v > fb.n {
			next <- v
			return
		}
		if v%5 != 0 && v%3 != 0 {
			if v == fb.n {
				fmt.Printf(" %d. ", v)
			} else {
				fmt.Printf(" %d,", v)
			}
			next <- v + 1
			continue
		}
		next <- v
	}
}
