package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	seatsLock sync.Mutex
	seats     int
	cond      = sync.NewCond(&seatsLock)
)

func main() {

}

func customers() {
	for {
		randomPause(1000)
	}
}

func randomPause(max int) {
	time.Sleep(time.Duration(rand.Intn(max)) * time.Millisecond)
}

func customer() {
	seatsLock.Lock()
	defer seatsLock.Unlock()
	if seats == 3 {
		log.Println("no free seats, the customer leaved")
		return
	}
	seats++
	cond.Broadcast()
	log.Println("customer start to barber")
}

func barber() {
	for {
		log.Println("tony wait for one customer to barber")
		seatsLock.Lock()
		for seats == 0 {
			cond.Wait()
		}
		seats--
		seatsLock.Unlock()
		log.Println("tony asked customer to barber")
		randomPause(2000)
	}
}
