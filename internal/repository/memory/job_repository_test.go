package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/repository/memory"
)

func TestJobRepositoryListNewestFirst(t *testing.T) {
	repo := memory.NewJobRepository()
	ctx := context.Background()
	older := time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)

	if err := repo.Save(ctx, &domain.Job{
		ID: "old", Type: domain.JobTypeBackup, Status: domain.JobStatusSucceeded, CreatedAt: older,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, &domain.Job{
		ID: "new", Type: domain.JobTypeBackup, Status: domain.JobStatusSucceeded, CreatedAt: newer,
	}); err != nil {
		t.Fatal(err)
	}

	list, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("len=%d", len(list))
	}
	if list[0].ID != "new" || list[1].ID != "old" {
		t.Fatalf("order=%s,%s want new,old", list[0].ID, list[1].ID)
	}
}
