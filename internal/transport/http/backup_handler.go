package http

import (
	"context"
	"io"
	"mime"
	nethttp "net/http"
	"strconv"

	"github.com/sboy99/go-vault/internal/domain"
)

func (r *router) handleListBackups(w nethttp.ResponseWriter, req *nethttp.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	status := req.URL.Query().Get("status")
	list, err := r.backup.ListBackups(req.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if status != "" {
		filtered := list[:0]
		for _, b := range list {
			if string(b.Status) == status {
				filtered = append(filtered, b)
			}
		}
		list = filtered
	}
	writeOK(w, list)
}

func (r *router) handleGetBackup(w nethttp.ResponseWriter, req *nethttp.Request) {
	id := req.PathValue("id")
	b, err := r.backup.GetBackup(req.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, b)
}

func (r *router) handleDownloadBackup(w nethttp.ResponseWriter, req *nethttp.Request) {
	id := req.PathValue("id")
	rc, b, err := r.backup.OpenBackup(req.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	defer func() { _ = rc.Close() }()
	w.Header().Set("Content-Type", "application/octet-stream")
	cd := mime.FormatMediaType("attachment", map[string]string{"filename": b.Name})
	if cd == "" {
		cd = "attachment"
	}
	w.Header().Set("Content-Disposition", cd)
	_, _ = io.Copy(w, rc)
}

func (r *router) handleCreateBackup(w nethttp.ResponseWriter, req *nethttp.Request) {
	j, err := r.jobs.TryStart(domain.JobTypeBackup, "", func(ctx context.Context) error {
		_, err := r.backup.BackupAndPrune(ctx)
		return err
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, j)
}
