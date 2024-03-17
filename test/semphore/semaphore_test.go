package semphore

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"runtime"
	"testing"
	"time"
)

//type semaphore struct {
//	sync.Locker
//	ch chan struct{}
//}
//
//func NewSemaphore(capacity int) sync.Locker {
//	if capacity <= 0 {
//		capacity = 1
//	}
//	return semaphore{
//		ch: make(chan struct{}, capacity),
//	}
//}
//
//func (s *semaphore) Lock() {
//	s.ch <- struct{}{}
//}
//
//func (s *semaphore) Unlock() {
//	<-s.ch
//}

var (
	maxWorkers = runtime.GOMAXPROCS(0)                    // worker amount
	sema       = semaphore.NewWeighted(int64(maxWorkers)) // semaphore
	task       = make([]int, maxWorkers*4)                // task amount, 4 times of the workers
)

func TestSemaphore(t *testing.T) {
	ctx := context.Background()

	for i := range task {
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}

		go func(i int) {
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond)
			task[i] = i + 1
		}(i)
	}
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		fmt.Printf("retrieve the whole worker failed, err: %v", err)
	}

	fmt.Println(task)

}
