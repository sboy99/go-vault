package api

import (
	"net/http"
	"strings"
)

func (s *Server) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := s.cfg.API.Token; token != "" {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != token {
				if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" {
					writeError(w, http.StatusUnauthorized, "unauthorized")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
