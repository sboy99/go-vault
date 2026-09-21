package api

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/sboy99/go-vault/internal/job"
)

func (s *Server) handleListBackups(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	status := r.URL.Query().Get("status")
	list, err := s.svc.ListBackups(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
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

func (s *Server) handleGetBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBackup(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeOK(w, b)
}

func (s *Server) handleDownloadBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rc, b, err := s.svc.OpenBackup(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	defer func() { _ = rc.Close() }()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+b.Name+"\"")
	_, _ = io.Copy(w, rc)
}

func (s *Server) handleCreateBackup(w http.ResponseWriter, r *http.Request) {
	j, err := s.runner.TryStart(job.TypeBackup, "", func(ctx context.Context) error {
		_, err := s.svc.BackupAndPrune(ctx)
		return err
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeOK(w, j)
}
