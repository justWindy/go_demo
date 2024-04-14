package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Chopstick struct {
	sync.Mutex
}

type Philosopher struct {
	id          int
	name        string
	left, right *Chopstick
	status      string
	sync.Mutex
}

func (p *Philosopher) dine() {
	for {
		mark(p, "meditation")
		randomPause(10)
		mark(p, "starve")
		p.left.Lock()
		mark(p, "retrieve the left chopstick")
		randomPause(100)
		p.right.Lock()
		mark(p, "eating")
		p.right.Unlock()
		p.left.Unlock()
	}
}

func randomPause(max int) {
	time.Sleep(time.Duration(rand.Intn(max)) * time.Millisecond)
}

func mark(p *Philosopher, action string) {
	fmt.Printf("%s start %s\n", p.name, action)
	p.status = fmt.Sprintf("%s start %s\n", p.status, action)
}

func main() {
	go http.ListenAndServe("localhost:8972", nil)

	count := 5
	chopsticks := make([]*Chopstick, count)
	for i := 0; i < count; i++ {
		chopsticks[i] = &Chopstick{}
	}

	names := []string{color.RedString("philosopher1"), color.MagentaString("philosopher2"),
		color.CyanString("philosopher3"), color.GreenString("philosopher4"), color.BlueString("philosopher5")}
	philosophers := make([]*Philosopher, count)
	for i := 0; i < count; i++ {
		philosophers[i] = &Philosopher{
			id:    i,
			name:  names[i],
			left:  chopsticks[i],
			right: chopsticks[(i+1)%count],
		}
		go philosophers[i].dine3()
	}
	//philosophers := minusDineCount()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("quit...Every philosopher status:")
	for _, p := range philosophers {
		fmt.Print(p.status)
	}
}

func minusDineCount() []*Philosopher {
	count := 5
	chopsticks := make([]*Chopstick, count)
	for i := 0; i < count; i++ {
		chopsticks[i] = &Chopstick{}
	}

	names := []string{color.RedString("philosopher1"), color.MagentaString("philosopher2"),
		color.CyanString("philosopher3"), color.GreenString("philosopher4"), color.BlueString("philosopher5")}
	philosophers := make([]*Philosopher, count)
	for i := 0; i < count; i++ {
		philosophers[i] = &Philosopher{
			name:  names[i],
			left:  chopsticks[i],
			right: chopsticks[(i+1)%count],
		}
		if i < count-1 {
			go philosophers[i].dine()
		}
	}

	return philosophers
}

func (p *Philosopher) dine1() {
	for {
		mark(p, "thinking")
		randomPause(10)
		mark(p, "starving")
		if p.id%2 == 1 {
			p.left.Lock()
			mark(p, "retrieve the left chopstick")
			p.right.Lock()
			mark(p, "eating")
			randomPause(10)
			p.right.Unlock()
			p.left.Unlock()
		} else {
			p.right.Lock()
			mark(p, "retrieve the right chopstick")
			p.left.Lock()
			mark(p, "eating")
			randomPause(10)
			p.left.Unlock()
			p.right.Unlock()
		}
	}
}

func (p *Philosopher) dine2() {
	for {
		mark(p, "thinking")
		randomPause(10)
		mark(p, "starving")
		if p.id == 4 {
			p.right.Lock()
			mark(p, "retrieve the right chopstick")
			p.left.Lock()
			mark(p, "eating")
			randomPause(10)
			p.left.Unlock()
			p.right.Unlock()
		} else {
			p.left.Lock()
			mark(p, "retrieve the left chopstick")
			p.right.Lock()
			mark(p, "eating")
			randomPause(10)
			p.right.Unlock()
			p.left.Unlock()
		}
	}
}

func (p *Philosopher) dine3() {
	for {
		mark(p, "thinking")
		randomPause(10)
		mark(p, "starving")
		p.Lock()
		p.left.Lock()
		mark(p, "retrieve the left chopstick")
		p.right.Lock()
		p.Unlock()
		mark(p, "eating")
		randomPause(10)
		p.right.Unlock()
		p.left.Unlock()
	}
}
