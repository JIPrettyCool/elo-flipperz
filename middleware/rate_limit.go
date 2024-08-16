package middleware

import (
    "golang.org/x/time/rate"
    "net/http"
    "sync"
)

var (
    limiter = rate.NewLimiter(1, 3) // 1 req per sec with 3 size
    mu      sync.Mutex
)

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mu.Lock()
        defer mu.Unlock()
        if !limiter.Allow() {
            http.Error(w, "Rate Limited biiaaachhh", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}