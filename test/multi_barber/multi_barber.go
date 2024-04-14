package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

var seats = make(Semaphore, 10)

func main() {
	go barber("Tony")
	go barber("Kevin")
	go barber("Allen")
	go customers()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func randomPause(max int) {
	time.Sleep(time.Duration(rand.Intn(max)) * time.Millisecond)
}

func barber(name string) {
	for {
		log.Println(name + " mr(miss) tony require for one customer")
		seats.Release()
		log.Println(name + " mr(miss) tony asked one customer, start barber")
		randomPause(2000)
	}
}

func customer() {
	if ok := seats.TryAcquire(); ok {
		log.Println("one customer sit down, and queue for barber")
	} else {
		log.Println("no free seat, one customer leaving")
	}
}

func customers() {
	for {
		randomPause(1000)
		go customer()
	}
}
