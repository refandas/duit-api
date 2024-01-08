package middleware

import (
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/web"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

// visitor represents an entity that visits or interacts with a system.
// It includes a rate limiter to control the frequency of interactions and
// keep track of the last time the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mutex sync.Mutex

func init() {
	go cleanupVisitors()
}

// getVisitor takes an IP address as a parameter and returns a rate.Limiter.
// The rate limiter enforces a constraint of allowing a maximum of 3 request
// per second for each visitor.
func getVisitor(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	guest, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 3)
		visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	guest.lastSeen = time.Now()
	return guest.limiter
}

// cleanupVisitors iterates through the visitors and removes any visitor whose
// last activity time exceeds 3 minutes.
func cleanupVisitors() {
	for {
		// Run cleanupVisitors() every 1 minute.
		time.Sleep(time.Minute)

		mutex.Lock()
		for ip, guest := range visitors {
			// Remove the visitors if there is no activity within 3 minutes.
			if time.Since(guest.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mutex.Unlock()
	}
}

// RateLimitMiddleware can be used to enforce rate limiting on incoming request
// to a web server, preventing abuse or excessive use.
type RateLimitMiddleware struct {
	Handler http.Handler
}

// NewRateLimitMiddleware takes an existing HTTP handler and returns a new
// RateLimitMiddleware instance, which can be used to enforce rate limiting
// on incoming requests to the provided handler.
func NewRateLimitMiddleware(handler http.Handler) *RateLimitMiddleware {
	return &RateLimitMiddleware{Handler: handler}
}

// ServeHTTP method satisfies the http.Handler interface and is responsible for
// handling incoming HTTP requests. It enforces rate limiting on requests based on
// the configured rate limit parameters.
func (middleware *RateLimitMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Get the IP address from the current visitor.
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		panic(err)
	}

	helper.SetupSecurityHeaders(writer)

	limiter := getVisitor(ip)
	if !limiter.Allow() {
		writer.WriteHeader(http.StatusTooManyRequests)

		webResponse := web.WebResponse{
			Code:   http.StatusTooManyRequests,
			Status: "TOO MANY REQUESTS",
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		middleware.Handler.ServeHTTP(writer, request)
	}
}
