package err_group

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestErrGroupUsing(t *testing.T) {
	var g errgroup.Group

	// start the first child goroutine
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		return nil
	})

	// start the second
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})

	// start the third
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		return nil
	})

	// wait for all the three goroutine done
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed: ", err)
	}
}
