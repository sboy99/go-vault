package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"

	"github.com/sboy99/go-vault/internal/domain"
)

type restoreRequest struct {
	BackupID string `json:"backup_id"`
	Confirm  string `json:"confirm"`
}

func (r *router) handleRestore(w nethttp.ResponseWriter, req *nethttp.Request) {
	var body restoreRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid json body")
		return
	}
	if body.BackupID == "" {
		writeError(w, nethttp.StatusBadRequest, "backup_id is required")
		return
	}
	if err := r.backup.ValidateRestoreConfirm(body.Confirm); err != nil {
		writeDomainError(w, err)
		return
	}
	j, err := r.jobs.TryStart(domain.JobTypeRestore, body.BackupID, func(ctx context.Context) error {
		return r.backup.RestoreBackup(ctx, body.BackupID, body.Confirm)
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, j)
}
