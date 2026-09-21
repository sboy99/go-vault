package api

import (
	"net/http"

	"github.com/sboy99/go-vault/internal/metrics"
)

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]string{"status": "ok"})
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	if !metrics.Ready() {
		writeError(w, http.StatusServiceUnavailable, "not ready")
		return
	}
	if err := s.svc.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	writeOK(w, map[string]string{"status": "ready"})
}
