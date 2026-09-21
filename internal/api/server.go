package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/backup"
	"github.com/sboy99/go-vault/internal/job"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/pkg/logger"
)

type Server struct {
	cfg    *config.Config
	svc    *backup.Service
	runner *job.Runner
	mux    *http.ServeMux
	http   *http.Server
}

type envelope struct {
	Success bool    `json:"success"`
	Data    any     `json:"data"`
	Error   *string `json:"error"`
}

func NewServer(cfg *config.Config, svc *backup.Service, runner *job.Runner) *Server {
	s := &Server{cfg: cfg, svc: svc, runner: runner, mux: http.NewServeMux()}
	s.routes()
	s.http = &http.Server{
		Addr:              cfg.API.Addr,
		Handler:           s.withMiddleware(s.mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /readyz", s.handleReadyz)
	s.mux.Handle("GET /metrics", metrics.Handler())
	s.mux.HandleFunc("GET /v1/backups", s.handleListBackups)
	s.mux.HandleFunc("GET /v1/backups/{id}", s.handleGetBackup)
	s.mux.HandleFunc("GET /v1/backups/{id}/download", s.handleDownloadBackup)
	s.mux.HandleFunc("POST /v1/backups", s.handleCreateBackup)
	s.mux.HandleFunc("POST /v1/restores", s.handleRestore)
	s.mux.HandleFunc("GET /v1/jobs", s.handleListJobs)
	s.mux.HandleFunc("GET /v1/jobs/{id}", s.handleGetJob)
}

// Handler returns the HTTP handler (useful for tests).
func (s *Server) Handler() http.Handler {
	return s.withMiddleware(s.mux)
}

func (s *Server) Start() error {
	logger.Info("api listening on %s", s.cfg.API.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

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
		_, err := s.svc.CreateBackup(ctx)
		if err != nil {
			return err
		}
		_, err = s.svc.Prune(ctx)
		return err
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeOK(w, j)
}

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

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := job.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeOK(w, j)
}

func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	jobs, err := job.List(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, jobs)
}

func writeOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(envelope{Success: true, Data: data, Error: nil})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(envelope{Success: false, Data: nil, Error: &msg})
}
