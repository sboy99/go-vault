package api

import (
	"context"
	"net/http"
	"time"

	"github.com/sboy99/go-vault/internal/backup"
	"github.com/sboy99/go-vault/internal/config"
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
