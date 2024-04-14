package main

import "log"

type Semaphore chan struct{}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}

func (s Semaphore) Release() {
	<-s
}

func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

var seats1 = make(Semaphore, 3)

func barber1() {
	for {
		log.Println("tony try to ask the on customer")
		seats1.Release()
		log.Println("tony asked one customer, and start to barber")
		randomPause(2000)
	}
}

func customer2() {
	if ok := seats1.TryAcquire(); ok {
		log.Println("one customer is sit down, wait for barber")
	} else {
		log.Println("no free seat, one customer leaving")
	}
}

func customers1() {
	for {
		randomPause(1000)
		go customer2()
	}
}
