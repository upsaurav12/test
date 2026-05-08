package middleware

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipRateLimiter maintains per-IP token-bucket rate limiters.
//
// NOTE: This implementation is in-process only. In a horizontally scaled
// deployment use a shared store (e.g. Redis) for distributed rate limiting.
// The limiters map grows unbounded; add a periodic cleanup mechanism if the
// set of client IPs is very large.
type ipRateLimiter struct {
	limiters sync.Map
	r        rate.Limit
	b        int
	ttl      time.Duration
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen int64
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	rl := &ipRateLimiter{
		r:   r,
		b:   b,
		ttl: 15 * time.Minute,
	}
	go rl.startCleanupLoop(5 * time.Minute)
	return rl
}

func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	now := time.Now().UnixNano()
	// Fast path: limiter already exists.
	if v, ok := i.limiters.Load(ip); ok {
		entry := v.(*limiterEntry)
		atomic.StoreInt64(&entry.lastSeen, now)
		return entry.limiter
	}
	// Slow path: create and store, guarding against a concurrent insert.
	entry := &limiterEntry{
		limiter:  rate.NewLimiter(i.r, i.b),
		lastSeen: now,
	}
	v, _ := i.limiters.LoadOrStore(ip, entry)
	stored := v.(*limiterEntry)
	atomic.StoreInt64(&stored.lastSeen, now)
	return stored.limiter
}

func (i *ipRateLimiter) startCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for now := range ticker.C {
		cutoff := now.Add(-i.ttl).UnixNano()
		i.limiters.Range(func(key, value any) bool {
			entry, ok := value.(*limiterEntry)
			if !ok || atomic.LoadInt64(&entry.lastSeen) < cutoff {
				i.limiters.Delete(key)
			}
			return true
		})
	}
}

// RateLimit returns a per-IP rate-limiting middleware using a token-bucket
// algorithm.
//   - r: sustained request rate (requests per second).
//   - b: burst size (maximum instantaneous requests allowed).
func RateLimit(r rate.Limit, b int) gin.HandlerFunc {
	limiter := newIPRateLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.getLimiter(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded – try again later",
			})
			return
		}
		c.Next()
	}
}
