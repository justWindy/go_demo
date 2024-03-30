package rate_limit

import (
	"github.com/juju/ratelimit"
	"log"
	"testing"
	"time"
)

func TestJujuLimiter(t *testing.T) {
	var bucket = ratelimit.NewBucket(time.Second, 3)
	for i := 0; i < 10; i++ {
		bucket.Wait(1)
		log.Printf("got #%d", i)
	}
}
