package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// fixedWindowIncrScript: increment counter; jika pertama kali, set expiry.
// Return: {count, ttl_ms}
const fixedWindowIncrScript = `
local key = KEYS[1]
local ttl_ms = tonumber(ARGV[1])
local count = redis.call("INCR", key)
if count == 1 then
	redis.call("PEXPIRE", key, ttl_ms)
end
local ttl = redis.call("PTTL", key)
return {count, ttl}
`

type fixedWindowResult struct {
	Count int64
	TTL   time.Duration
}

// incrFixedWindow menaikkan counter & memastikan key punya expiry.
func incrFixedWindow(ctx context.Context, client *redis.Client, key string, window time.Duration) (fixedWindowResult, error) {
	vals, err := client.Eval(ctx, fixedWindowIncrScript, []string{key}, window.Milliseconds()).Result()
	if err != nil {
		return fixedWindowResult{}, err
	}

	arr, ok := vals.([]interface{})
	if !ok || len(arr) < 2 {
		return fixedWindowResult{}, errUnexpectedRedisReply
	}

	count := toInt64(arr[0])
	ttl := toInt64(arr[1])

	return fixedWindowResult{
		Count: count,
		TTL:   time.Duration(ttl) * time.Millisecond,
	}, nil
}

func toInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	case string:
		var out int64
		for _, r := range n {
			if r < '0' || r > '9' {
				break
			}
			out = out*10 + int64(r-'0')
		}
		return out
	default:
		return 0
	}
}
