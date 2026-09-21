package http

import nethttp "net/http"

func (r *router) handleGetStats(w nethttp.ResponseWriter, req *nethttp.Request) {
	stats, err := r.backup.Stats(req.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, map[string]any{
		"backups":     stats,
		"job_running": r.jobs.Busy(),
	})
}

func (r *router) handleGetRetention(w nethttp.ResponseWriter, req *nethttp.Request) {
	preview, err := r.backup.RetentionPreview(req.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, preview)
}
