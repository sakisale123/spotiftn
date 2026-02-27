package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	requests int
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

const (
	maxRequests = 20
	window      = time.Minute
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		mu.Lock()
		c, exists := clients[ip]

		if !exists || time.Since(c.lastSeen) > window {
			clients[ip] = &client{
				requests: 1,
				lastSeen: time.Now(),
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if c.requests >= maxRequests {
			mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		c.requests++
		c.lastSeen = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
