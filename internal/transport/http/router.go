package http

import (
	"time"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/metrics"
	nethttp "net/http"
)

// RouterDeps holds services the HTTP transport needs.
type RouterDeps struct {
	Backup   *app.BackupService
	Jobs     *app.JobService
	Config   *config.Config
	NextRun  func() time.Time
	APIToken string
}

type router struct {
	backup  *app.BackupService
	jobs    *app.JobService
	cfg     *config.Config
	nextRun func() time.Time
}

// NewRouter registers API routes and returns an http.Handler.
// Server lifecycle (ListenAndServe / Shutdown) belongs in cmd/server.
func NewRouter(deps RouterDeps) nethttp.Handler {
	r := &router{
		backup:  deps.Backup,
		jobs:    deps.Jobs,
		cfg:     deps.Config,
		nextRun: deps.NextRun,
	}
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /healthz", r.handleHealthz)
	mux.HandleFunc("GET /readyz", r.handleReadyz)
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /v1/backups", r.handleListBackups)
	mux.HandleFunc("GET /v1/backups/{id}", r.handleGetBackup)
	mux.HandleFunc("GET /v1/backups/{id}/download", r.handleDownloadBackup)
	mux.HandleFunc("POST /v1/backups", r.handleCreateBackup)
	mux.HandleFunc("POST /v1/restores", r.handleRestore)
	mux.HandleFunc("GET /v1/jobs", r.handleListJobs)
	mux.HandleFunc("GET /v1/jobs/{id}", r.handleGetJob)
	mux.HandleFunc("GET /v1/config", r.handleGetConfig)
	mux.HandleFunc("GET /v1/stats", r.handleGetStats)
	mux.HandleFunc("GET /v1/retention", r.handleGetRetention)
	return withBearerAuth(deps.APIToken, mux)
}
