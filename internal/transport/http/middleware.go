package http

import (
	"crypto/subtle"
	nethttp "net/http"
	"strings"
)

func withBearerAuth(token string, next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if token != "" {
			auth := r.Header.Get("Authorization")
			provided := ""
			if strings.HasPrefix(auth, "Bearer ") {
				provided = strings.TrimPrefix(auth, "Bearer ")
			}
			if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
				if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" {
					writeError(w, nethttp.StatusUnauthorized, "unauthorized")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
