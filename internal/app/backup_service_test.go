package app_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/repository/memory"
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
	return domain.ObjectInfo{Key: key, Size: int64(len(data)), SHA256: "abc"}, nil
}

func (s *stubStore) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := s.objects[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return io.NopCloser(newBytesReader(data)), nil
}

func (s *stubStore) Delete(ctx context.Context, key string) error {
	delete(s.objects, key)
	return nil
}

func (s *stubStore) List(ctx context.Context, prefix string) ([]domain.ObjectInfo, error) {
	var out []domain.ObjectInfo
	for k, v := range s.objects {
		out = append(out, domain.ObjectInfo{Key: k, Size: int64(len(v))})
	}
	return out, nil
}

type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader { return &bytesReader{data: data} }

func (b *bytesReader) Read(p []byte) (int, error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

type stubEngine struct{}

func (stubEngine) DumpTo(ctx context.Context, p domain.ConnParams, w io.Writer) (*domain.DumpResult, error) {
	_, err := w.Write([]byte("dump-data"))
	return &domain.DumpResult{ServerVersion: "16.0", ServerMajor: 16}, err
}
func (stubEngine) Verify(ctx context.Context, dumpPath string, serverMajor int) error { return nil }
func (stubEngine) Restore(ctx context.Context, p domain.ConnParams, dumpPath string, jobs int) error {
	return nil
}
func (stubEngine) ServerVersion(ctx context.Context, p domain.ConnParams) (string, int, error) {
	return "16.0", 16, nil
}

type stubNotifier struct{ events []string }

func (n *stubNotifier) NotifyFailure(ctx context.Context, event, message string) {
	n.events = append(n.events, event)
}

type stubMetrics struct {
	failed  int
	success int
	restore []string
	pruned  int
}

func (m *stubMetrics) BackupSucceeded(started time.Time, size int64) { m.success++ }
func (m *stubMetrics) BackupFailed()                                 { m.failed++ }
func (m *stubMetrics) RestoreFinished(status string)                 { m.restore = append(m.restore, status) }
func (m *stubMetrics) BackupPruned()                                 { m.pruned++ }

func newTestBackupService(t *testing.T) (*app.BackupService, *memory.BackupRepository, *stubMetrics) {
	t.Helper()
	repo := memory.NewBackupRepository()
	metrics := &stubMetrics{}
	svc := app.NewBackupService(
		app.BackupConfig{
			DBName:         "appdb",
			DatabaseType:   domain.DatabasePostgreSQL,
			StorageType:    domain.StorageLocal,
			Conn:           domain.ConnParams{Host: "localhost", Port: 5432, Name: "appdb", Username: "u", Password: "p", SSLMode: "disable"},
			TempDir:        t.TempDir(),
			DumpTimeout:    time.Minute,
			RestoreTimeout: time.Minute,
			RestoreJobs:    1,
			Retention:      domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12},
			Location:       time.UTC,
		},
		repo,
		newStubStore(),
		stubEngine{},
		&stubNotifier{},
		metrics,
	)
	return svc, repo, metrics
}

func TestRestoreRequiresConfirm(t *testing.T) {
	svc, _, _ := newTestBackupService(t)
	err := svc.RestoreBackup(context.Background(), "any", "wrong")
	if !errors.Is(err, domain.ErrRestoreConfirmMismatch) {
		t.Fatalf("want ErrRestoreConfirmMismatch, got %v", err)
	}
}

func TestCreateBackupPersistsSuccess(t *testing.T) {
	svc, repo, metrics := newTestBackupService(t)
	b, err := svc.CreateBackup(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != domain.StatusSuccess {
		t.Fatalf("status=%s", b.Status)
	}
	got, err := repo.FindByID(context.Background(), b.BackupId)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusSuccess {
		t.Fatalf("persisted status=%s", got.Status)
	}
	if metrics.success != 1 {
		t.Fatalf("success metric=%d", metrics.success)
	}
}

func TestListBackups(t *testing.T) {
	svc, _, _ := newTestBackupService(t)
	if _, err := svc.CreateBackup(context.Background()); err != nil {
		t.Fatal(err)
	}
	list, err := svc.ListBackups(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list=%d", len(list))
	}
}
