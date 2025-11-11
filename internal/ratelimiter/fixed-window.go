package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	window  time.Duration
	limit   int
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		window:  window,
		limit:   limit,
	}
}

func (f *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	f.RLock()
	count, exists := f.clients[ip]
	f.RUnlock()
	if !exists || count <= f.limit {
		f.RLock()
		if !exists {
			go f.reset(ip)
		}
		f.clients[ip]++
		f.RUnlock()
		return true, 0
	}
	return false, f.window
}

// it will run after every window unit of time so that the users data from map after that time set to zero
func (f *FixedWindowRateLimiter) reset(ip string) {
	time.Sleep(f.window)
	f.RLock()
	delete(f.clients, ip)
	f.RUnlock()
}
