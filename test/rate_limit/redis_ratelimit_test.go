package rate_limit

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func TestRedisRateLimit(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		i := i
		go func() {
			defer wg.Done()

			ctx := context.Background()
			rdb := redis.NewClient(&redis.Options{
				Addr: "192.168.50.89:6379",
			})
			_ = rdb.FlushDB(ctx).Err()

			limiter := redis_rate.NewLimiter(rdb)
			for j := 0; j < 10; j++ {
				res, err := limiter.Allow(ctx, "token:123", redis_rate.PerSecond(5))
				if err != nil {
					panic(err)
				}
				log.Println(i, "allowed", res.Allowed, "remaining", res.Remaining, "retry after", res.RetryAfter)
				if res.Allowed == 0 {
					time.Sleep(res.RetryAfter)
				}
			}
		}()
	}

	wg.Wait()
}
