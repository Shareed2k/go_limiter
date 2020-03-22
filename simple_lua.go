package go_limiter

import "github.com/go-redis/redis/v7"

var script2 = redis.NewScript(`
-- this script has side-effects, so it requires replicate commands mode
redis.replicate_commands()

local rate_limit_key = KEYS[1]
local rate = tonumber(ARGV[1])
local period = tonumber(ARGV[2])

local now = redis.call("TIME")

local jan_1_2017 = 1483228800
local now_nanos = (now[1] - jan_1_2017) + (now[2] / 1000000)

local clear_before = now_nanos - period

local function allow_check_card ()
    redis.call("ZREMRANGEBYSCORE", rate_limit_key, "0.0", clear_before)

    return redis.call("ZCARD", rate_limit_key)
end

local function delta()
    local res = redis.call("ZRANGEBYSCORE", rate_limit_key, "0.0", now_nanos, "WITHSCORES", "limit", 0, 1)
    local oldest = 0

    if #res > 0 then
        oldest = res[2]
    end

    local gab = now_nanos - oldest

    return period - gab
end

local del = delta()
local count = allow_check_card()

if count >= rate then
    return {0, count, tostring(del)}
end

redis.call("ZADD", rate_limit_key, now_nanos, now_nanos)
redis.call("EXPIRE", rate_limit_key, period)

return {1, count+1, tostring(del)}
`)
