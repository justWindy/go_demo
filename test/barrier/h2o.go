package main

import (
	"context"
	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"time"
)

type H2O struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
	b     cyclicbarrier.CyclicBarrier
}

func New() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		b:     cyclicbarrier.New(3),
	}
}

func (h *H2O) hydrogen(releaseHydrogen func()) {
	h.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()
	h.b.Await(context.Background())
	h.semaH.Release(1)
}

func (h *H2O) oxygen(releaseOxygen func()) {
	h.semaO.Acquire(context.Background(), 1)
	releaseOxygen()
	h.b.Await(context.Background())
	h.semaO.Release(1)
}

func randomPause(max int) {
	time.Sleep(time.Duration(rand.Intn(max)) * time.Millisecond)
}
