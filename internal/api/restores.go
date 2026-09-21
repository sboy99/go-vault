package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sboy99/go-vault/internal/job"
)

type restoreRequest struct {
	BackupID string `json:"backup_id"`
	Confirm  string `json:"confirm"`
}

func (s *Server) handleRestore(w http.ResponseWriter, r *http.Request) {
	var req restoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if req.BackupID == "" {
		writeError(w, http.StatusBadRequest, "backup_id is required")
		return
	}
	if req.Confirm != s.cfg.DB.Name {
		writeError(w, http.StatusBadRequest, "confirm must match the configured database name")
		return
	}
	j, err := s.runner.TryStart(job.TypeRestore, req.BackupID, func(ctx context.Context) error {
		return s.svc.RestoreBackup(ctx, req.BackupID)
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeOK(w, j)
}
