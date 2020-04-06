// https://github.com/go-redis/redis_rate
package go_limiter

import (
	"strconv"
	"time"
)

const GCRAAlgorithmName = "gcra"

type gcra struct {
	key   string
	limit *Limit
	rdb   rediser
}

// Allow is shorthand for AllowN(key, 1).
func (c *gcra) Allow() (*Result, error) {
	return c.AllowN(1)
}

// SetKey _
func (c *gcra) SetKey(key string) {
	c.key = key
}

// AllowN reports whether n events may happen at time now.
func (c *gcra) AllowN(n int) (*Result, error) {
	limit := c.limit
	values := []interface{}{limit.Burst, limit.Rate, limit.Period.Seconds(), n}

	v, err := script.Run(c.rdb, []string{c.key}, values...).Result()
	if err != nil {
		return nil, err
	}

	values = v.([]interface{})

	retryAfter, err := strconv.ParseFloat(values[2].(string), 64)
	if err != nil {
		return nil, err
	}

	resetAfter, err := strconv.ParseFloat(values[3].(string), 64)
	if err != nil {
		return nil, err
	}

	res := &Result{
		Limit:      limit,
		Key:        c.key,
		Allowed:    values[0].(int64) == 0,
		Remaining:  values[1].(int64),
		RetryAfter: dur(retryAfter),
		ResetAfter: dur(resetAfter),
	}
	return res, nil
}

func dur(f float64) time.Duration {
	if f == -1 {
		return -1
	}
	return time.Duration(f * float64(time.Second))
}
