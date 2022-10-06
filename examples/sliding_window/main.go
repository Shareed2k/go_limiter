package main

import (
	"context"
	"github.com/go-redis/redis/v8"
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

	limiter := go_limiter.NewLimiter(client)
	res, err := limiter.Allow(context.Background(), "api_gateway:klu4ik", &go_limiter.Limit{
		Algorithm: go_limiter.SlidingWindowAlgorithm,
		Rate:      10,
		Period:    2 * time.Minute,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("===> ", res.Allowed, res.Remaining)
	// Output: true 1
}
