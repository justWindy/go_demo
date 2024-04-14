package go_break

import (
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func TestAvoidExecute(t *testing.T) {
	var flag atomic.Bool

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			for i := 0; i < 100; i++ {
				if !flag.CompareAndSwap(false, true) {
					time.Sleep(time.Second)
					continue
				}

				flag.Store(false)
			}
		}()
	}
	wg.Wait()
}
