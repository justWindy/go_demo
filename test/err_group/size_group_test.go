package err_group

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-pkgz/syncs"
	"github.com/mdlayher/schedgroup"
	"github.com/vardius/gollback"
)

func TestSizeGroupUsing(t *testing.T) {
	swg := syncs.NewSizedGroup(10)
	// swg = syncs.NewSizedGroup(10, syncs.Preemptive)
	var c uint32

	for i := 0; i < 1000; i++ {
		swg.Go(func(ctx context.Context) {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
		})
	}

	swg.Wait()
	fmt.Println(c)
}

func TestGollbackAll(t *testing.T) {
	rs, errs := gollback.All(
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed")
		}, func(ctx context.Context) (interface{}, error) {
			return 3, nil
		})

	fmt.Println(rs)
	fmt.Println(errs)
}

func TestGollBackRetry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		fmt.Println("execute times")
		return nil, errors.New("failed")
	})

	fmt.Println(res)
	fmt.Println(err)
}

func TestSchedGroup(t *testing.T) {
	sg := schedgroup.New(context.Background())

	for i := 0; i < 3; i++ {
		n := i + 1
		sg.Delay(time.Duration(n)*100*time.Millisecond, func() {
			fmt.Println("task:", n)
		})
	}

	if err := sg.Wait(); err != nil {
		fmt.Printf("failed to wait: %v", err)
		panic(err)
	}
}
