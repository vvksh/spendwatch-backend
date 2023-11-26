package main

import (
	"golang.org/x/time/rate"
	"net/http"
)

func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(2, 4)
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("The API is at capacity, try again later."))
			return
		} else {
			next(w, r)
		}
	}
}
