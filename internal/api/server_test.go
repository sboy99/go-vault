package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/api"
	"github.com/sboy99/go-vault/internal/backup"
	"github.com/sboy99/go-vault/internal/job"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/internal/storage"
	"github.com/spf13/viper"
)

func TestHealthz(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestListBackupsEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/backups", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
	var env map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env["success"] != true {
		t.Fatalf("envelope: %v", env)
	}
}

func TestRestoreRequiresConfirm(t *testing.T) {
	h := newTestHandler(t)
	body := bytes.NewBufferString(`{"backup_id":"x","confirm":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/restores", body)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBearerAuthWhenConfigured(t *testing.T) {
	dir := t.TempDir()
	_ = setupTestConfig(t, dir)
	viper.Set("api.token", "secret")
	cfg := config.GetConfig()

	store, err := storage.NewStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	svc := backup.NewService(cfg, store)
	runner := job.NewRunner()
	s := api.NewServer(cfg, svc, runner)
	h := s.Handler()

	req := httptest.NewRequest(http.MethodGet, "/v1/backups", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/v1/backups", nil)
	req2.Header.Set("Authorization", "Bearer secret")
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr2.Code, rr2.Body.String())
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	dir := t.TempDir()
	cfg := setupTestConfig(t, dir)
	store, err := storage.NewStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	svc := backup.NewService(cfg, store)
	runner := job.NewRunner()
	metrics.SetReady(true)
	return api.NewServer(cfg, svc, runner).Handler()
}

func setupTestConfig(t *testing.T, dir string) *config.Config {
	t.Helper()
	viper.Set("app.name", "go-vault")
	viper.Set("app.version", "test")
	viper.Set("db.name", "app")
	viper.Set("db.type", "POSTGRESQL")
	viper.Set("db.host", "localhost")
	viper.Set("db.port", 5432)
	viper.Set("db.username", "postgres")
	viper.Set("db.password", "postgres")
	viper.Set("db.sslmode", "disable")
	viper.Set("storage.type", "LOCAL")
	viper.Set("storage.dest", filepath.Join(dir, "backups"))
	viper.Set("schedule.cron", "0 2 * * *")
	viper.Set("schedule.timezone", "UTC")
	viper.Set("retention.daily", 7)
	viper.Set("retention.weekly", 4)
	viper.Set("retention.monthly", 12)
	viper.Set("api.addr", ":0")
	viper.Set("api.token", "")
	viper.Set("runtime.temp_dir", filepath.Join(dir, "tmp"))
	viper.Set("runtime.dump_timeout", time.Hour)
	viper.Set("runtime.restore_timeout", time.Hour)
	viper.Set("runtime.shutdown_timeout", time.Minute)
	viper.Set("runtime.meta_db_path", filepath.Join(dir, "meta.db"))
	viper.Set("runtime.restore_jobs", 2)

	metaPath := filepath.Join(dir, "meta.db")
	_ = meta.Cleanup()
	if err := meta.Init(metaPath); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = meta.Cleanup() })

	cfg := config.GetConfig()
	_ = os.MkdirAll(cfg.Storage.Dest, 0o755)
	_ = os.MkdirAll(cfg.Runtime.TempDir, 0o755)
	return cfg
}
