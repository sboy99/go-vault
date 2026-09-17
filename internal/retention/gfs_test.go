package retention_test

import (
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/internal/retention"
)

func TestSelectKeepsNewestAlways(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, loc)
	backups := []meta.BackupMeta{
		backup("a", now.Add(-1*time.Hour)),
		backup("b", now.Add(-48*time.Hour)),
	}
	keep, prune := retention.Select(backups, retention.Policy{Daily: 1, Weekly: 1, Monthly: 1}, now, loc)
	if len(keep) < 1 {
		t.Fatal("expected at least newest kept")
	}
	if keep[0].BackupId != "a" {
		t.Fatalf("expected newest a, got %s", keep[0].BackupId)
	}
	_ = prune
}

func TestSelectGFSDailyWeeklyMonthly(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 3, 31, 15, 0, 0, 0, loc)

	var backups []meta.BackupMeta
	// 40 daily backups going back
	for i := 0; i < 40; i++ {
		ts := now.AddDate(0, 0, -i).Add(-time.Hour)
		backups = append(backups, backup(ts.Format("2006-01-02"), ts))
	}

	keep, prune := retention.Select(backups, retention.Policy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	if len(keep) == 0 {
		t.Fatal("keep empty")
	}
	if len(keep)+len(prune) != len(backups) {
		t.Fatalf("keep+prune=%d want %d", len(keep)+len(prune), len(backups))
	}

	ids := map[string]struct{}{}
	for _, b := range keep {
		ids[b.BackupId] = struct{}{}
	}
	// Newest day must be kept
	if _, ok := ids[now.AddDate(0, 0, 0).Add(-time.Hour).Format("2006-01-02")]; !ok {
		// id is the day key we used
		newestID := now.Add(-time.Hour).Format("2006-01-02")
		if _, ok := ids[newestID]; !ok {
			t.Fatalf("newest day not kept: %v", ids)
		}
	}
	// Should prune more than half of 40 with 7/4/12
	if len(prune) < 10 {
		t.Fatalf("expected significant prune, got keep=%d prune=%d", len(keep), len(prune))
	}
}

func TestSelectIgnoresFailed(t *testing.T) {
	loc := time.UTC
	now := time.Now().UTC()
	backups := []meta.BackupMeta{
		{BackupId: "ok", Status: meta.StatusSuccess, CreatedAt: now},
		{BackupId: "bad", Status: meta.StatusFailed, CreatedAt: now.Add(-time.Hour)},
	}
	keep, prune := retention.Select(backups, retention.Policy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	if len(keep) != 1 || keep[0].BackupId != "ok" {
		t.Fatalf("keep=%v", keep)
	}
	if len(prune) != 0 {
		t.Fatalf("failed backups should not be pruned via GFS, got %v", prune)
	}
}

func TestSelectMonthBoundary(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 1, 5, 12, 0, 0, 0, loc)
	backups := []meta.BackupMeta{
		backup("jan", time.Date(2024, 1, 5, 1, 0, 0, 0, loc)),
		backup("dec", time.Date(2023, 12, 15, 1, 0, 0, 0, loc)),
		backup("old", time.Date(2022, 6, 1, 1, 0, 0, 0, loc)),
	}
	keep, prune := retention.Select(backups, retention.Policy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	ids := map[string]bool{}
	for _, b := range keep {
		ids[b.BackupId] = true
	}
	if !ids["jan"] || !ids["dec"] {
		t.Fatalf("expected jan and dec kept, keep=%v", keep)
	}
	if !containsID(prune, "old") {
		t.Fatalf("expected old pruned, prune=%v", prune)
	}
}

func TestSelectDSTSafe(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip(err)
	}
	// Around US DST spring forward 2024-03-10
	now := time.Date(2024, 3, 12, 12, 0, 0, 0, loc)
	var backups []meta.BackupMeta
	for i := 0; i < 10; i++ {
		ts := now.AddDate(0, 0, -i)
		backups = append(backups, backup(ts.Format("2006-01-02"), ts))
	}
	keep, _ := retention.Select(backups, retention.Policy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	if len(keep) < 7 {
		t.Fatalf("expected at least 7 daily keeps around DST, got %d", len(keep))
	}
}

func backup(id string, at time.Time) meta.BackupMeta {
	return meta.BackupMeta{
		BackupId:  id,
		Name:      id + ".dump",
		Status:    meta.StatusSuccess,
		CreatedAt: at,
	}
}

func containsID(list []meta.BackupMeta, id string) bool {
	for _, b := range list {
		if b.BackupId == id {
			return true
		}
	}
	return false
}
