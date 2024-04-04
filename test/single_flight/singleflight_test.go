package single_flight

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"

	"golang.org/x/sync/singleflight"
)

var (
	cache = make(map[string]string)
	sf    singleflight.Group
	mu    sync.Mutex
)

func fetchData(key string) (string, error) {
	fmt.Println("fetching data for key:", key)
	return "Data for " + key, nil
}

func getDataWithCache(key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if data, ok := cache[key]; ok {
		return data, nil
	}

	v, err, _ := sf.Do(key, func() (interface{}, error) {
		return fetchData(key)
	})

	if err != nil {
		return "", err
	}

	cache[key] = v.(string)
	return v.(string), nil
}

func TestGetDataWithCache(t *testing.T) {
	var wg sync.WaitGroup

	keys := []string{"key1", "key2", "key1", "key1", "key1"}

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			data, err := getDataWithCache(k)
			if err != nil {
				fmt.Println("Error fetching data for key: ", k, "-", err)
			} else {
				fmt.Println("Data for key", k, ":", data)
			}
		}(key)
	}
	wg.Wait()
}

var errGoexit = errors.New("runtime.Goexit was called")

type result struct {
	Val    interface{}
	Err    error
	Shared bool
}

type panicError struct {
	value interface{}
	stack []byte
}

func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error

	forgotten bool
	dups      int
	chans     []chan<- result
}

type group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	return c.val, c.err, c.dups > 0
}

func (g *group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	defer func() {
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				go panic(e)
				select {}
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
		} else {
			for _, ch := range c.chans {
				ch <- result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()
		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}
