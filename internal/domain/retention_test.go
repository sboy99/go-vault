package domain_test

import (
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
)

func policy() domain.RetentionPolicy {
	return domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12}
}

func TestPlanOpenWeekKeepsAllDailies(t *testing.T) {
	loc := time.UTC
	// Wednesday 2024-09-11; week Sun 9/8 – Sat 9/14 is still open.
	now := time.Date(2024, 9, 11, 15, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("mon", time.Date(2024, 9, 9, 2, 0, 0, 0, loc)),
		backup("tue", time.Date(2024, 9, 10, 2, 0, 0, 0, loc)),
		backup("wed", time.Date(2024, 9, 11, 2, 0, 0, 0, loc)),
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if len(plan.Keep) != 3 || len(plan.Prune) != 0 {
		t.Fatalf("keep=%d prune=%d want 3/0", len(plan.Keep), len(plan.Prune))
	}
	for _, b := range plan.Keep {
		if plan.Tiers[b.BackupId] != domain.TierDaily {
			t.Fatalf("%s tier=%s want daily", b.BackupId, plan.Tiers[b.BackupId])
		}
	}
}

func TestPlanClosedWeekKeepsOnlySaturday(t *testing.T) {
	loc := time.UTC
	// Monday 2024-09-16; prior week Sun 9/8 – Sat 9/14 is closed.
	now := time.Date(2024, 9, 16, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("sun", time.Date(2024, 9, 8, 2, 0, 0, 0, loc)),
		backup("wed", time.Date(2024, 9, 11, 2, 0, 0, 0, loc)),
		backup("sat", time.Date(2024, 9, 14, 2, 0, 0, 0, loc)),
		backup("mon", time.Date(2024, 9, 16, 2, 0, 0, 0, loc)), // open week
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if !containsID(plan.Keep, "sat") {
		t.Fatalf("expected sat weekly kept, keep=%v", idsOf(plan.Keep))
	}
	if plan.Tiers["sat"] != domain.TierWeekly {
		t.Fatalf("sat tier=%s want weekly", plan.Tiers["sat"])
	}
	if containsID(plan.Keep, "sun") || containsID(plan.Keep, "wed") {
		t.Fatalf("closed-week non-anchors should be pruned, keep=%v", idsOf(plan.Keep))
	}
	if !containsID(plan.Keep, "mon") || plan.Tiers["mon"] != domain.TierDaily {
		t.Fatalf("open-week mon should be daily kept")
	}
}

func TestPlanMissingSaturdayFallsBackToNewest(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 9, 16, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("sun", time.Date(2024, 9, 8, 2, 0, 0, 0, loc)),
		backup("fri", time.Date(2024, 9, 13, 2, 0, 0, 0, loc)), // newest in week, no Sat
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if !containsID(plan.Keep, "fri") {
		t.Fatalf("expected fri fallback weekly, keep=%v", idsOf(plan.Keep))
	}
	if plan.Tiers["fri"] != domain.TierWeekly {
		t.Fatalf("fri tier=%s want weekly", plan.Tiers["fri"])
	}
	if containsID(plan.Keep, "sun") {
		t.Fatalf("sun should be pruned")
	}
}

func TestPlanClosedMonthKeepsMonthEndDropsWeeklies(t *testing.T) {
	loc := time.UTC
	// 2024-10-05; September is closed.
	now := time.Date(2024, 10, 5, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("sep-sat", time.Date(2024, 9, 28, 2, 0, 0, 0, loc)), // Saturday
		backup("sep-30", time.Date(2024, 9, 30, 2, 0, 0, 0, loc)),  // month end (Monday)
		backup("oct-fri", time.Date(2024, 10, 4, 2, 0, 0, 0, loc)),
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if !containsID(plan.Keep, "sep-30") || plan.Tiers["sep-30"] != domain.TierMonthly {
		t.Fatalf("sep-30 should be monthly, tiers=%v keep=%v", plan.Tiers, idsOf(plan.Keep))
	}
	if containsID(plan.Keep, "sep-sat") {
		t.Fatalf("sep weekly should be pruned after month close, keep=%v", idsOf(plan.Keep))
	}
	if !containsID(plan.Keep, "oct-fri") {
		t.Fatalf("oct open-week daily should be kept")
	}
}

func TestPlanMissingMonthEndFallsBack(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 10, 5, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("sep-20", time.Date(2024, 9, 20, 2, 0, 0, 0, loc)),
		backup("sep-28", time.Date(2024, 9, 28, 2, 0, 0, 0, loc)), // newest in Sep, no 30th
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if !containsID(plan.Keep, "sep-28") || plan.Tiers["sep-28"] != domain.TierMonthly {
		t.Fatalf("sep-28 should be monthly fallback, keep=%v tiers=%v", idsOf(plan.Keep), plan.Tiers)
	}
	if containsID(plan.Keep, "sep-20") {
		t.Fatalf("sep-20 should be pruned")
	}
}

func TestPlanWeekStraddlingMonthBoundary(t *testing.T) {
	loc := time.UTC
	// Week Sun 2024-09-29 – Sat 2024-10-05. Now is Oct 7 so week and September are closed.
	now := time.Date(2024, 10, 7, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("sep-29", time.Date(2024, 9, 29, 2, 0, 0, 0, loc)), // Sunday
		backup("sep-30", time.Date(2024, 9, 30, 2, 0, 0, 0, loc)), // month end
		backup("oct-1", time.Date(2024, 10, 1, 2, 0, 0, 0, loc)),
		backup("oct-5", time.Date(2024, 10, 5, 2, 0, 0, 0, loc)), // Saturday week anchor
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if plan.Tiers["sep-30"] != domain.TierMonthly {
		t.Fatalf("sep-30 tier=%s want monthly", plan.Tiers["sep-30"])
	}
	if containsID(plan.Keep, "sep-29") {
		t.Fatalf("sep-29 pruned by closed month")
	}
	if plan.Tiers["oct-5"] != domain.TierWeekly {
		t.Fatalf("oct-5 tier=%s want weekly (closed week in open month)", plan.Tiers["oct-5"])
	}
	if containsID(plan.Keep, "oct-1") {
		t.Fatalf("oct-1 should be pruned as non-anchor of closed week")
	}
}

func TestPlanFailedNeverPruned(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 9, 16, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("ok", time.Date(2024, 9, 14, 2, 0, 0, 0, loc)),
		{BackupId: "bad", Status: domain.StatusFailed, CreatedAt: time.Date(2024, 9, 13, 2, 0, 0, 0, loc)},
	}
	plan := domain.Plan(backups, policy(), now, loc)
	if !containsID(plan.Keep, "bad") {
		t.Fatalf("failed should be kept, keep=%v", idsOf(plan.Keep))
	}
	if plan.Tiers["bad"] != domain.TierDaily {
		t.Fatalf("failed tier=%s want daily", plan.Tiers["bad"])
	}
	if containsID(plan.Prune, "bad") {
		t.Fatalf("failed must not be pruned")
	}
}

func TestPlanCapsTrimExcess(t *testing.T) {
	loc := time.UTC
	// Open month with many closed weeks of Saturdays.
	now := time.Date(2024, 9, 25, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("w1", time.Date(2024, 8, 31, 2, 0, 0, 0, loc)), // Sat
		backup("w2", time.Date(2024, 9, 7, 2, 0, 0, 0, loc)),  // Sat
		backup("w3", time.Date(2024, 9, 14, 2, 0, 0, 0, loc)), // Sat
		backup("w4", time.Date(2024, 9, 21, 2, 0, 0, 0, loc)), // Sat
		backup("today", time.Date(2024, 9, 25, 2, 0, 0, 0, loc)),
	}
	// Weekly cap 2 → only two newest weeklies (w4, w3); w2/w1 pruned by cap.
	// August closed → w1 is monthly for Aug (Aug 31 is month end), not weekly.
	plan := domain.Plan(backups, domain.RetentionPolicy{Daily: 7, Weekly: 2, Monthly: 12}, now, loc)
	if plan.Tiers["w1"] != domain.TierMonthly {
		t.Fatalf("w1 (Aug 31) should be monthly, got %s", plan.Tiers["w1"])
	}
	weeklyKept := 0
	for _, b := range plan.Keep {
		if plan.Tiers[b.BackupId] == domain.TierWeekly {
			weeklyKept++
		}
	}
	if weeklyKept > 2 {
		t.Fatalf("weekly cap 2, got %d weeklies in keep", weeklyKept)
	}
	if !containsID(plan.Keep, "today") {
		t.Fatalf("newest daily must be kept")
	}
}

func TestPlanNewestAlwaysKept(t *testing.T) {
	loc := time.UTC
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, loc)
	backups := []domain.Backup{
		backup("a", now.Add(-1*time.Hour)),
		backup("b", now.Add(-48*time.Hour)),
	}
	plan := domain.Plan(backups, domain.RetentionPolicy{Daily: 1, Weekly: 1, Monthly: 1}, now, loc)
	if len(plan.Keep) < 1 || plan.Keep[0].BackupId != "a" {
		t.Fatalf("expected newest a first, keep=%v", idsOf(plan.Keep))
	}
}

func TestPlanTimezoneBoundary(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	// 2024-10-01 03:30 UTC = 2024-09-30 23:30 EDT → still Sep 30 in NY.
	utcInstant := time.Date(2024, 10, 1, 3, 30, 0, 0, time.UTC)
	now := time.Date(2024, 10, 5, 12, 0, 0, 0, ny)
	backups := []domain.Backup{
		backup("sep-end", utcInstant),
		backup("oct-fri", time.Date(2024, 10, 4, 12, 0, 0, 0, ny)),
	}
	plan := domain.Plan(backups, policy(), now, ny)
	if plan.Tiers["sep-end"] != domain.TierMonthly {
		t.Fatalf("sep-end tier=%s want monthly (Sep 30 in NY)", plan.Tiers["sep-end"])
	}
}

func backup(id string, at time.Time) domain.Backup {
	return domain.Backup{
		BackupId:  id,
		Name:      id + ".dump",
		Status:    domain.StatusSuccess,
		CreatedAt: at,
	}
}

func containsID(list []domain.Backup, id string) bool {
	for _, b := range list {
		if b.BackupId == id {
			return true
		}
	}
	return false
}

func idsOf(list []domain.Backup) []string {
	ids := make([]string, len(list))
	for i, b := range list {
		ids[i] = b.BackupId
	}
	return ids
}
