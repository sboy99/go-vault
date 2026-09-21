package http

import (
	nethttp "net/http"
	"strings"
)

func withBearerAuth(token string, next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if token != "" {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != token {
				if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" {
					writeError(w, nethttp.StatusUnauthorized, "unauthorized")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
