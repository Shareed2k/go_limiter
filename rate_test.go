package go_limiter

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

func rateLimiter() *Limiter {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	if err := client.FlushDB().Err(); err != nil {
		panic(err)
	}
	return NewLimiter(client)
}

func TestLimiter_Allow(t *testing.T) {
	l := rateLimiter()

	limit := &Limit{
		Algorithm: SlidingWindowAlgorithm,
		Rate:      10,
		Period:    time.Minute,
		Burst:     10,
	}

	t.Run("simple", func(t *testing.T) {
		res, err := l.Allow("test_me", limit)

		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, int64(9), res.Remaining)
	})

	t.Run("gcra", func(t *testing.T) {
		limit.Algorithm = GCRAAlgorithm

		res, err := l.Allow("test_me", limit)

		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, int64(9), res.Remaining)
		assert.Equal(t, res.RetryAfter, time.Duration(-1))
	})
}

func BenchmarkAllow_Simple(b *testing.B) {
	l := rateLimiter()
	limit := &Limit{
		Algorithm: SlidingWindowAlgorithm,
		Rate:      10000,
		Period:    time.Second,
		Burst:     10000,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := l.Allow("foo", limit)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllow_GCRA(b *testing.B) {
	l := rateLimiter()
	limit := &Limit{
		Algorithm: SlidingWindowAlgorithm,
		Rate:      10000,
		Period:    time.Second,
		Burst:     10000,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := l.Allow("foo", limit)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
