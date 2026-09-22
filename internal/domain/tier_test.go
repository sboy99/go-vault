package domain_test

import (
	"testing"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
)

func TestPlanCycleHelpersViaWeekAndMonth(t *testing.T) {
	loc := time.UTC
	// Closed week only: Sat promoted.
	now := time.Date(2024, 9, 16, 12, 0, 0, 0, loc)
	plan := domain.Plan([]domain.Backup{
		backup("sat", time.Date(2024, 9, 14, 2, 0, 0, 0, loc)),
		backup("sun", time.Date(2024, 9, 8, 2, 0, 0, 0, loc)),
	}, domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	if plan.Tiers["sat"] != domain.TierWeekly {
		t.Fatalf("Saturday of closed week want weekly, got %s", plan.Tiers["sat"])
	}

	// Closed month: last day promoted even when not Saturday.
	now = time.Date(2024, 10, 2, 12, 0, 0, 0, loc)
	plan = domain.Plan([]domain.Backup{
		backup("sep-30", time.Date(2024, 9, 30, 2, 0, 0, 0, loc)), // Monday
		backup("sep-28", time.Date(2024, 9, 28, 2, 0, 0, 0, loc)), // Saturday
	}, domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12}, now, loc)
	if plan.Tiers["sep-30"] != domain.TierMonthly {
		t.Fatalf("month-end want monthly, got %s", plan.Tiers["sep-30"])
	}
	if containsID(plan.Keep, "sep-28") {
		t.Fatalf("prior Saturday should be pruned after month close")
	}
}

func TestPlanNilLocationDefaultsUTC(t *testing.T) {
	now := time.Date(2024, 10, 2, 12, 0, 0, 0, time.UTC)
	plan := domain.Plan([]domain.Backup{
		backup("sep-30", time.Date(2024, 9, 30, 2, 0, 0, 0, time.UTC)),
	}, domain.RetentionPolicy{Daily: 7, Weekly: 4, Monthly: 12}, now, nil)
	if plan.Tiers["sep-30"] != domain.TierMonthly {
		t.Fatalf("got %s want monthly with nil loc", plan.Tiers["sep-30"])
	}
}
