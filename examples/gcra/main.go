package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	"github.com/shareed2k/go_limiter"
)

func main() {
	option, err := redis.ParseURL("redis://127.0.0.1:6379/0")
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(option)
	ctx := context.Background()

	limiter := go_limiter.NewLimiter(client)
	res, err := limiter.Allow(ctx, "api_gateway:klu4ik", &go_limiter.Limit{
		Algorithm: go_limiter.GCRAAlgorithm,
		Rate:      10,
		Period:    2 * time.Minute,
		Burst:     10,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("===> ", res.Allowed, res.Remaining)
	// Output: true 1
}
