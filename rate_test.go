package go_limiter

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"

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

	if err := client.FlushDB(context.Background()).Err(); err != nil {
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

	ctx := context.Background()

	t.Run("simple", func(t *testing.T) {
		res, err := l.Allow(ctx, "test_me", limit)

		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, int64(9), res.Remaining)
	})

	t.Run("gcra", func(t *testing.T) {
		limit.Algorithm = GCRAAlgorithm

		res, err := l.Allow(ctx, "test_me", limit)

		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, int64(9), res.Remaining)
		assert.Equal(t, res.RetryAfter, time.Duration(-1))
	})
}

func TestLimiter_Reset(t *testing.T) {
	l := rateLimiter()

	limit := &Limit{
		Rate:   1,
		Period: time.Minute,
		Burst:  1,
	}

	ctx := context.Background()

	t.Run("reset-sliding-window", func(t *testing.T) {
		limit.Algorithm = SlidingWindowAlgorithm

		res, err := l.Allow(ctx, "sliding-reset_me", limit)
		assert.Nil(t, err)
		assert.True(t, res.Allowed)

		res, err = l.Allow(ctx, "sliding-reset_me", limit)
		assert.Nil(t, err)
		assert.False(t, res.Allowed)

		err = l.Reset(ctx, "sliding-reset_me", limit)
		assert.Nil(t, err)

		res, err = l.Allow(ctx, "sliding-reset_me", limit)
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
	})

	t.Run("reset-gcra-window", func(t *testing.T) {
		limit.Algorithm = GCRAAlgorithm

		res, err := l.Allow(ctx, "gcra-reset_me", limit)
		assert.Nil(t, err)
		assert.True(t, res.Allowed)

		res, err = l.Allow(ctx, "gcra-reset_me", limit)
		assert.Nil(t, err)
		assert.False(t, res.Allowed)

		err = l.Reset(ctx, "gcra-reset_me", limit)
		assert.Nil(t, err)

		res, err = l.Allow(ctx, "gcra-reset_me", limit)
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
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

	ctx := context.Background()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := l.Allow(ctx, "foo", limit)
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

	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := l.Allow(ctx, "foo", limit)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
