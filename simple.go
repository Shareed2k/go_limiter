package go_limiter

import (
	"strconv"
)

type simple struct {
	key   string
	limit *Limit
	rdb   rediser
}

func (c *simple) Allow() (r *Result, err error) {
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
		Allowed:    values[0].(int64) == 1,
		Remaining:  values[1].(int64),
		RetryAfter: dur(retryAfter),
		ResetAfter: limit.Period,
	}, nil
}
