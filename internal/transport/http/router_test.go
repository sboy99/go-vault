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
	"github.com/sboy99/go-vault/internal/config"
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

	cfg := &config.Config{
		App: config.App{Name: "go-vault", Version: "test"},
		DB: config.Database{
			Name:     "app",
			Type:     config.POSTGRESQL,
			Host:     "localhost",
			Port:     5432,
			Username: "u",
			Password: "super-secret-password",
			SSLMode:  "disable",
		},
		Storage: config.Storage{
			Type: config.LOCAL,
			Dest: "/tmp/backups",
			Cloud: config.CloudStorage{
				Type: config.AWS,
				AWS: config.AWSCloudStorage{
					Region:          "us-east-1",
					BucketName:      "vault",
					AccessKeyId:     "AKIA...",
					AccessKeySecret: "aws-secret-should-not-leak",
					Endpoint:        "",
				},
			},
		},
		Schedule:  config.Schedule{Cron: "0 2 * * *", Timezone: "UTC"},
		Retention: config.Retention{Daily: 7, Weekly: 4, Monthly: 12},
		API:       config.API{Addr: ":8080", Token: "api-token-should-not-leak"},
	}

	return transporthttp.NewRouter(transporthttp.RouterDeps{
		Backup:   backupSvc,
		Jobs:     jobSvc,
		Config:   cfg,
		NextRun:  func() time.Time { return time.Date(2026, 9, 23, 2, 0, 0, 0, time.UTC) },
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

func TestGetConfigRedactsSecrets(t *testing.T) {
	h := newTestHandler(t, "")
	req := httptest.NewRequest(nethttp.MethodGet, "/v1/config", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusOK {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, secret := range []string{
		"super-secret-password",
		"aws-secret-should-not-leak",
		"api-token-should-not-leak",
		"password",
		"access_key_secret",
		"token",
	} {
		if bytes.Contains(rr.Body.Bytes(), []byte(secret)) && secret != "password" && secret != "token" && secret != "access_key_secret" {
			t.Fatalf("response leaked secret %q: %s", secret, body)
		}
	}
	// Field names for secrets must not appear either.
	for _, field := range []string{"\"password\"", "\"access_key_secret\"", "\"token\"", "\"access_key_id\""} {
		if bytes.Contains(rr.Body.Bytes(), []byte(field)) {
			t.Fatalf("response included secret field %s: %s", field, body)
		}
	}
	var env struct {
		Success bool `json:"success"`
		Data    struct {
			DB struct {
				Name string `json:"name"`
			} `json:"db"`
			Schedule struct {
				Cron string `json:"cron"`
			} `json:"schedule"`
			Retention struct {
				Daily int `json:"daily"`
			} `json:"retention"`
			NextRun *time.Time `json:"next_run"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if !env.Success || env.Data.DB.Name != "app" || env.Data.Schedule.Cron != "0 2 * * *" || env.Data.Retention.Daily != 7 {
		t.Fatalf("unexpected config payload: %+v", env.Data)
	}
	if env.Data.NextRun == nil {
		t.Fatal("expected next_run")
	}
}

func TestGetStatsAndRetention(t *testing.T) {
	h := newTestHandler(t, "")

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/stats", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != nethttp.StatusOK {
		t.Fatalf("stats status %d body=%s", rr.Code, rr.Body.String())
	}
	var statsEnv map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &statsEnv); err != nil {
		t.Fatal(err)
	}
	if statsEnv["success"] != true {
		t.Fatalf("stats envelope: %v", statsEnv)
	}
	data, ok := statsEnv["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing data: %v", statsEnv)
	}
	if _, ok := data["job_running"]; !ok {
		t.Fatalf("missing job_running: %v", data)
	}
	if _, ok := data["backups"]; !ok {
		t.Fatalf("missing backups: %v", data)
	}

	req2 := httptest.NewRequest(nethttp.MethodGet, "/v1/retention", nil)
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)
	if rr2.Code != nethttp.StatusOK {
		t.Fatalf("retention status %d body=%s", rr2.Code, rr2.Body.String())
	}
}
