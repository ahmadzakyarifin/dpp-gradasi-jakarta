package middleware

import (
	"sync"
	"time"
)

// memoryLimiter adalah implementasi fixed-window rate limit in-memory.
// Aman untuk single-instance (gaya DPP murni, tanpa Redis).
// Map key → counter + window expiry; cleanup otomatis saat key kedaluwarsa.
type memoryLimiter struct {
	mu     sync.Mutex
	prefix string
	items  map[string]*windowCounter
}

type windowCounter struct {
	count int64
	until time.Time
}

func newMemoryLimiter() *memoryLimiter {
	return &memoryLimiter{
		prefix: "dppgradasi_rate_limit",
		items:  make(map[string]*windowCounter),
	}
}

// check menaikkan counter; jika melebihi limit → blocked=true + retryAfter.
func (l *memoryLimiter) check(key string, limit int64, window time.Duration) (blocked bool, retryAfter int, remaining int64) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	fullKey := l.prefix + ":" + key

	item, ok := l.items[fullKey]
	if !ok || now.After(item.until) {
		// Window baru (atau yang lama sudah kedaluwarsa).
		item = &windowCounter{
			count: 0,
			until: now.Add(window),
		}
		l.items[fullKey] = item
	}

	item.count++

	remaining = limit - item.count
	if remaining < 0 {
		remaining = 0
	}

	retryAfter = int(time.Until(item.until).Seconds())
	if retryAfter < 1 {
		retryAfter = 1
	}

	return item.count > limit, retryAfter, remaining
}
