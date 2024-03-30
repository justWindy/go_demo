package rate_limit

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// 令牌产生速率，每200ms产生一个令牌
	var limit = rate.Every(200 * time.Millisecond)

	var limiter = rate.NewLimiter(limit, 3) // 令牌桶的容量为3

	for i := 0; i < 15; i++ {
		log.Printf("got #%d, err: %v\n", i, limiter.Wait(context.Background()))
	}
}

func TestSetLimitAt(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var limit = rate.Every(200 * time.Millisecond)
	var limiter = rate.NewLimiter(limit, 3)

	for i := 0; i < 3; i++ {
		log.Printf("got #%d, err: %v", i, limiter.Wait(context.Background()))
	}

	log.Println("set new limit as 10s")
	limiter.SetLimitAt(time.Now().Add(10*time.Second), rate.Every(3*time.Second))

	for i := 4; i < 9; i++ {
		log.Printf("got #%d, err: %v", i, limiter.Wait(context.Background()))
	}
}

func TestReserveN(t *testing.T) {
	var limiter = rate.NewLimiter(1, 10)
	limiter.WaitN(context.Background(), 10) // 把初始的令牌消耗掉

	r := limiter.ReserveN(time.Now().Add(5), 4)
	log.Printf("ok : %v, delay: %v", r.OK(), r.Delay())
	r.Cancel()

	r = limiter.ReserveN(time.Now().Add(3), 6)
	log.Printf("ok: %v, delay: %v", r.OK(), r.Delay())

	r = limiter.ReserveN(time.Now().Add(3), 100)
	log.Printf("ok: %v, delay: %t", r.OK(), r.Delay() == rate.InfDuration)
}
