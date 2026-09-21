package http

import (
	"github.com/sboy99/go-vault/internal/metrics"
	nethttp "net/http"
)

func (r *router) handleHealthz(w nethttp.ResponseWriter, req *nethttp.Request) {
	writeOK(w, map[string]string{"status": "ok"})
}

func (r *router) handleReadyz(w nethttp.ResponseWriter, req *nethttp.Request) {
	if !metrics.Ready() {
		writeError(w, nethttp.StatusServiceUnavailable, "not ready")
		return
	}
	if err := r.backup.Ping(req.Context()); err != nil {
		writeError(w, nethttp.StatusServiceUnavailable, err.Error())
		return
	}
	writeOK(w, map[string]string{"status": "ready"})
}
