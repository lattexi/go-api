package middleware

import (
	"fmt"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example middleware logic
		w.Header().Set("X-Test-Middleware", "Success")
		fmt.Println("TestMiddleware executed")
		next.ServeHTTP(w, r)
	})
}
