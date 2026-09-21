package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/internal/repository/memory"
	transporthttp "github.com/sboy99/go-vault/internal/transport/http"
)

type stubStore struct {
	objects map[string][]byte
}

func newStubStore() *stubStore {
	return &stubStore{objects: map[string][]byte{}}
}

func (s *stubStore) Save(ctx context.Context, key string, r io.Reader) (domain.ObjectInfo, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return domain.ObjectInfo{}, err
	}
	s.objects[key] = data
	return domain.ObjectInfo{Key: key, Size: int64(len(data))}, nil
}
func (s *stubStore) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(s.objects[key])), nil
}
func (s *stubStore) Delete(ctx context.Context, key string) error { delete(s.objects, key); return nil }
func (s *stubStore) List(ctx context.Context, prefix string) ([]domain.ObjectInfo, error) {
	return nil, nil
}

type stubEngine struct{}

func (stubEngine) DumpTo(ctx context.Context, p domain.ConnParams, w io.Writer) (*domain.DumpResult, error) {
	_, _ = w.Write([]byte("x"))
	return &domain.DumpResult{ServerMajor: 16}, nil
}
func (stubEngine) Verify(ctx context.Context, dumpPath string, serverMajor int) error { return nil }
func (stubEngine) Restore(ctx context.Context, p domain.ConnParams, dumpPath string, jobs int) error {
	return nil
}
func (stubEngine) ServerVersion(ctx context.Context, p domain.ConnParams) (string, int, error) {
	return "16", 16, nil
}

type noopNotifier struct{}

func (noopNotifier) NotifyFailure(ctx context.Context, event, message string) {}

type noopMetrics struct{}

func (noopMetrics) BackupSucceeded(started time.Time, size int64) {}
func (noopMetrics) BackupFailed()                                 {}
func (noopMetrics) RestoreFinished(status string)                 {}
func (noopMetrics) BackupPruned()                                 {}

func newTestHandler(t *testing.T, token string) nethttp.Handler {
	t.Helper()
	backupRepo := memory.NewBackupRepository()
	jobRepo := memory.NewJobRepository()
	backupSvc := app.NewBackupService(
		app.BackupConfig{
			DBName:         "app",
			DatabaseType:   domain.DatabasePostgreSQL,
			StorageType:    domain.StorageLocal,
			Conn:           domain.ConnParams{Host: "localhost", Port: 5432, Name: "app", Username: "u", Password: "p", SSLMode: "disable"},
			TempDir:        t.TempDir(),
			DumpTimeout:    time.Minute,
			RestoreTimeout: time.Minute,
			RestoreJobs:    1,
			Retention:      domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12},
			Location:       time.UTC,
		},
		backupRepo,
		newStubStore(),
		stubEngine{},
		noopNotifier{},
		noopMetrics{},
	)
	jobSvc := app.NewJobService(jobRepo)
	metrics.SetReady(true)
	return transporthttp.NewRouter(transporthttp.RouterDeps{
		Backup:   backupSvc,
		Jobs:     jobSvc,
		APIToken: token,
	})
}

func TestHealthz(t *testing.T) {
	h := newTestHandler(t, "")
	req := httptest.NewRequest(nethttp.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestListBackupsEmpty(t *testing.T) {
	h := newTestHandler(t, "")
	req := httptest.NewRequest(nethttp.MethodGet, "/v1/backups", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusOK {
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
	h := newTestHandler(t, "")
	body := bytes.NewBufferString(`{"backup_id":"x","confirm":"wrong"}`)
	req := httptest.NewRequest(nethttp.MethodPost, "/v1/restores", body)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusBadRequest {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRestoreMissingBackupID(t *testing.T) {
	h := newTestHandler(t, "")
	body := bytes.NewBufferString(`{"backup_id":"","confirm":"app"}`)
	req := httptest.NewRequest(nethttp.MethodPost, "/v1/restores", body)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusBadRequest {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBearerAuthWhenConfigured(t *testing.T) {
	h := newTestHandler(t, "secret")

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/backups", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	req2 := httptest.NewRequest(nethttp.MethodGet, "/v1/backups", nil)
	req2.Header.Set("Authorization", "Bearer secret")
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)
	if rr2.Code != nethttp.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr2.Code, rr2.Body.String())
	}
}
