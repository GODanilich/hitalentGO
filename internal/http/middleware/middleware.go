package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"GODanilich/hitalentGO/internal/apperr"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					writeJSON(w, http.StatusInternalServerError, map[string]any{
						"error": map[string]any{
							"code":    apperr.CodeInternal,
							"message": "panic recovered",
						},
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s in %s", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
