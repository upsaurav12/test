package middleware

import (
	"net/http"
	"sync"

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
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	return &ipRateLimiter{r: r, b: b}
}

func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	// Fast path: limiter already exists.
	if v, ok := i.limiters.Load(ip); ok {
		return v.(*rate.Limiter)
	}
	// Slow path: create and store, guarding against a concurrent insert.
	l := rate.NewLimiter(i.r, i.b)
	v, _ := i.limiters.LoadOrStore(ip, l)
	return v.(*rate.Limiter)
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
