package middlewares

import (
	"net/http"
	"sync"
	"time"
)

func (rl *RateLimiter) RateLimiters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		ip := r.RemoteAddr
		rl.visitors[ip]++
		if rl.visitors[ip] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type RateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	rl.resetVisitorCount()
	return rl
}

func (rl *RateLimiter) resetVisitorCount() {
	go func() {
		for {
			time.Sleep(rl.resetTime)
			rl.mu.Lock()
			rl.visitors = make(map[string]int)
			rl.mu.Unlock()
		}
	}()
}
