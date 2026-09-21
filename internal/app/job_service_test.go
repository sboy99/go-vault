package app_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/repository/memory"
)

func TestJobServiceSingleFlight(t *testing.T) {
	jobs := memory.NewJobRepository()
	svc := app.NewJobService(jobs)

	started := make(chan struct{})
	release := make(chan struct{})

	j1, err := svc.TryStart(domain.JobTypeBackup, "", func(ctx context.Context) error {
		close(started)
		<-release
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	<-started

	_, err = svc.TryStart(domain.JobTypeRestore, "x", func(ctx context.Context) error { return nil })
	if !errors.Is(err, domain.ErrJobAlreadyRunning) {
		t.Fatalf("want ErrJobAlreadyRunning, got %v", err)
	}

	close(release)
	deadline := time.Now().Add(2 * time.Second)
	for svc.Busy() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if svc.Busy() {
		t.Fatal("still busy")
	}

	got, err := svc.Get(context.Background(), j1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.JobStatusSucceeded {
		t.Fatalf("status=%s", got.Status)
	}
}

func TestJobServiceList(t *testing.T) {
	jobs := memory.NewJobRepository()
	svc := app.NewJobService(jobs)

	var wg sync.WaitGroup
	wg.Add(1)
	_, err := svc.TryStart(domain.JobTypeBackup, "", func(ctx context.Context) error {
		defer wg.Done()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	list, err := svc.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list=%d", len(list))
	}
}

func TestSetupServiceRequiresFields(t *testing.T) {
	store := &settingsStore{}
	svc := app.NewSetupService(store)
	err := svc.SaveSettings(domain.Settings{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSetupServiceSaves(t *testing.T) {
	store := &settingsStore{}
	svc := app.NewSetupService(store)
	err := svc.SaveSettings(domain.Settings{
		DBName:      "app",
		DBHost:      "localhost",
		DBUsername:  "postgres",
		StorageType: domain.StorageLocal,
	})
	if err != nil {
		t.Fatal(err)
	}
	if store.saved.DBName != "app" {
		t.Fatalf("saved=%+v", store.saved)
	}
}

type settingsStore struct {
	saved domain.Settings
}

func (s *settingsStore) Save(settings domain.Settings) error {
	s.saved = settings
	return nil
}
