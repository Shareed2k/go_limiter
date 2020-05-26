package go_limiter

import (
	"strconv"
)

const SlidingWindowAlgorithmName = "sliding_window"

type slidingWindow struct {
	key   string
	limit *Limit
	rdb   rediser
}

func (c *slidingWindow) Reset() error {
	res := c.rdb.Del(c.key)
	return res.Err()
}

func (c *slidingWindow) SetKey(key string) {
	c.key = key
}

func (c *slidingWindow) Allow() (r *Result, err error) {
	limit := c.limit
	values := []interface{}{limit.Rate, limit.Period.Seconds()}

	v, err := script2.Run(c.rdb, []string{c.key}, values...).Result()
	if err != nil {
		return nil, err
	}

	values = v.([]interface{})

	retryAfter, err := strconv.ParseFloat(values[2].(string), 64)
	if err != nil {
		return nil, err
	}

	return &Result{
		Limit:      limit,
		Key:        c.key,
		Allowed:    values[0].(int64) == 1,
		Remaining:  values[1].(int64),
		RetryAfter: dur(retryAfter),
		ResetAfter: limit.Period,
	}, nil
}
