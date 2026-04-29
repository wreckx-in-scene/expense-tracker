package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// each IP get its own limiter
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

// function 1 : get client
func getClient(ip string) *rate.Limiter {
	//locking the map first
	mu.Lock()
	defer mu.Unlock()

	//checking if this ip has a limiter
	if c, ok := clients[ip]; ok {
		c.lastSeen = time.Now()
		return c.limiter
	}

	// if IP doesnt exist --> create a new limiter
	// 20 requests per minute , burst of 20
	limiter := rate.NewLimiter(rate.Every(time.Minute/20), 20)

	//add to map
	clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

func CleanupClients() {
	//locking the mutex
	mu.Lock()
	defer mu.Unlock()

	//searching for the ip
	for ip, c := range clients {
		if time.Since(c.lastSeen) > 3*time.Minute {
			delete(clients, ip)
		}
	}
}

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//getting ip address from request
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		limiter := getClient(ip)

		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}

		next(w, r)
	}
}
