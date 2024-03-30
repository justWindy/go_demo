package rate_limit

import (
	"go.uber.org/ratelimit"
	"log"
	"testing"
	"time"
)

func TestUberRateLimit(t *testing.T) {
	rl := ratelimit.New(1, ratelimit.WithSlack(3))

	for i := 0; i < 10; i++ {
		rl.Take()
		log.Printf("got #%d", i)
		if i == 3 {
			time.Sleep(5 * time.Second)
		}
	}
}