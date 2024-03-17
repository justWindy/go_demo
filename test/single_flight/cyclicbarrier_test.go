package single_flight

import (
	"context"
	"errors"
	"fmt"
	"github.com/marusama/cyclicbarrier"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type cyclicBarry2 interface {
	//等待所有的参与者到达，如果被ctx.Done()中断，则会返回ErrBrokenBarrier
	Await(ctx context.Context) error

	// 重置循环屏障到初始状态。如果当前有等待着，那么它们会返回ErrBrokenBarrier
	Reset()

	// 返回当前等待者的数量
	GetNumberWaiting() int

	// 参与者的数量
	GetParties() int

	// 循环屏障是否处于中断状态
	IsBroken() bool
}

type round struct {
	count    int           // 这一轮参与的goroutine的数量
	waitCh   chan struct{} // 这一轮等待channel
	brokeCh  chan struct{} // 广播用的channel
	isBroken bool          // 屏障是否认为破坏
}

type cyclicBarry struct {
	parties       int          // 参与者的数量
	barrierAction func() error // 屏障打开时要调用的函数

	lock  sync.RWMutex
	round *round // 轮次
}

func (b *cyclicBarry) Await(ctx context.Context) error {
	var (
		ctxDoneCh <-chan struct{}
	)

	if ctx != nil {
		ctxDoneCh = ctx.Done()
	}

	// 检查ctx是否已经被取消或者超时
	select {
	case <-ctxDoneCh:
		return ctx.Err()
	default:
	}
	//加锁
	b.lock.Lock()

	// 如果这一轮的等待和释放已经完成
	if b.round.isBroken {
		b.lock.Unlock()
		return errors.New("errBrokenBarrier")
	}

	// 在这一轮数据将调用的参与者数量加1
	b.round.count++

	// 先保存这一轮的相关对象备用，避免发生数据竞争，获取新一轮的对象
	waitCh := b.round.waitCh
	brokeCh := b.round.brokeCh
	count := b.round.count

	b.lock.Unlock()

	//下面就不需要锁了，因为本轮的对象已经获取到本地变量了
	if count > b.parties {
		panic("CyclicBarrier.Await is called more than count of parties")
	}

	// 如果当前的调用者不是最后一个调用者，则被阻塞等待
	if count < b.parties {
		// 等待发生下面的情况之一
		// 1. 最后一个调用者到来
		// 2. 人为破坏本轮的等待
		// 3. ctx被完成
		select {
		case <-waitCh:
			return nil
		case <-brokeCh:
			return errors.New("errBrokenBarrier")
		case <-ctxDoneCh:
			return ctx.Err()
		}
	} else {
		if b.barrierAction != nil {
			err := b.barrierAction()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func TestCyclicBarrier(t *testing.T) {
	cnt := 0
	b := cyclicbarrier.NewWithAction(10, func() error {
		cnt++
		return nil
	})

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				fmt.Printf("%v goroutine %d 来到第%d轮屏障\n", time.Now(), i, j)
				err := b.Await(context.TODO())
				fmt.Printf("%v goroutine %d 冲破第%d轮屏障\n", time.Now(), i, j)
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println(cnt)
}
