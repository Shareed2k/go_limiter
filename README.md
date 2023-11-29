# Rate limiting with few algorithms (Sliding Window, Leaky Bucket)

[![Build Status](https://travis-ci.org/Shareed2k/go_limiter.svg?branch=master)](github.com/shareed2k/go_limiter)

This package is based on [go-redis/redis_rate](https://github.com/go-redis/redis_rate) and implements GCRA (aka leaky bucket) for rate limiting based on Redis. The code requires Redis version 3.2 or newer since it relies on replicate_commands feature.

## Installation

go_limiter requires a Go version with [Modules](https://github.com/golang/go/wiki/Modules) support and uses import versioning. So please make sure to initialize a Go module before installing go_limiter:

```shell
go get github.com/shareed2k/go_limiter
```

Import:
```go
import "github.com/shareed2k/go_limiter"
```

## Examplle
```go
import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/shareed2k/go_limiter"
)

func main() {
	option, err := redis.ParseURL("redis://127.0.0.1:6379/0")
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(option)
	_ = client.FlushDB().Err()

	limiter := go_limiter.NewLimiter(client)
	res, err := limiter.Allow("api_gateway_cache:klu4ik", &go_limiter.Limit{
		// or you can use go_limiter.SlidingWindowAlgorithm
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
 
```
