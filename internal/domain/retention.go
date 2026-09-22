package domain

import (
	"math"
	"sort"
	"time"
)

// RetentionPolicy caps how many backups of each rollup tier to keep.
type RetentionPolicy struct {
	Daily   int
	Weekly  int
	Monthly int
}

// RetentionPlan is the keep/prune split with per-backup rollup tiers.
type RetentionPlan struct {
	Tiers map[string]BackupTier
	Keep  []Backup
	Prune []Backup
}

// Plan classifies successful backups by cycle rollup and returns keep/prune.
//
// Week cycle is Sunday–Saturday; month cycle is the calendar month. Closing a
// week promotes its Saturday backup (or newest in the week) to weekly and
// prunes the rest. Closing a month promotes its last-calendar-day backup (or
// newest in the month) to monthly and prunes the rest of that month.
// Non-successful backups are never pruned and are labeled daily. Caps trim the
// newest N of each tier. The newest successful backup is always kept.
func Plan(backups []Backup, p RetentionPolicy, now time.Time, loc *time.Location) RetentionPlan {
	if loc == nil {
		loc = time.UTC
	}
	now = now.In(loc)

	out := RetentionPlan{
		Tiers: map[string]BackupTier{},
		Keep:  []Backup{},
		Prune: []Backup{},
	}

	var success []Backup
	var other []Backup
	for _, b := range backups {
		if b.Status == StatusSuccess {
			success = append(success, b)
		} else {
			other = append(other, b)
		}
	}

	for _, b := range other {
		b.Tier = TierDaily
		out.Tiers[b.BackupId] = TierDaily
		out.Keep = append(out.Keep, b)
	}

	if len(success) == 0 {
		sortKeepPrune(&out)
		return out
	}

	sort.Slice(success, func(i, j int) bool {
		return success[i].CreatedAt.After(success[j].CreatedAt)
	})
	newestID := success[0].BackupId

	byWeek := map[string][]Backup{}
	byMonth := map[string][]Backup{}
	for _, b := range success {
		wk := weekKey(b.CreatedAt, loc)
		mk := monthKey(b.CreatedAt, loc)
		byWeek[wk] = append(byWeek[wk], b)
		byMonth[mk] = append(byMonth[mk], b)
	}

	weekAnchors := map[string]string{}
	for wk, group := range byWeek {
		end := weekCycleEnd(group[0].CreatedAt, loc)
		if !cycleClosed(end, now, loc) {
			continue
		}
		weekAnchors[wk] = resolveAnchor(group, loc, func(t time.Time) bool {
			return dateOnly(t, loc).Equal(dateOnly(end, loc))
		})
	}

	monthAnchors := map[string]string{}
	for mk, group := range byMonth {
		end := monthCycleEnd(group[0].CreatedAt, loc)
		if !cycleClosed(end, now, loc) {
			continue
		}
		monthAnchors[mk] = resolveAnchor(group, loc, func(t time.Time) bool {
			return dateOnly(t, loc).Equal(dateOnly(end, loc))
		})
	}

	type classified struct {
		b    Backup
		tier BackupTier
		keep bool
	}
	var candidates []classified

	for _, b := range success {
		mk := monthKey(b.CreatedAt, loc)
		wk := weekKey(b.CreatedAt, loc)
		monthEnd := monthCycleEnd(b.CreatedAt, loc)
		weekEnd := weekCycleEnd(b.CreatedAt, loc)

		var tier BackupTier
		keep := false

		switch {
		case cycleClosed(monthEnd, now, loc):
			if monthAnchors[mk] == b.BackupId {
				tier = TierMonthly
				keep = true
			}
		case cycleClosed(weekEnd, now, loc):
			if weekAnchors[wk] == b.BackupId {
				tier = TierWeekly
				keep = true
			}
		default:
			tier = TierDaily
			keep = true
		}

		if keep {
			candidates = append(candidates, classified{b: b, tier: tier, keep: true})
		} else {
			b.Tier = TierDaily
			out.Tiers[b.BackupId] = TierDaily
			out.Prune = append(out.Prune, b)
		}
	}

	// Apply safety caps per tier (newest N). Newest successful always survives.
	byTier := map[BackupTier][]classified{
		TierDaily:   {},
		TierWeekly:  {},
		TierMonthly: {},
	}
	for _, c := range candidates {
		byTier[c.tier] = append(byTier[c.tier], c)
	}
	for tier, list := range byTier {
		sort.Slice(list, func(i, j int) bool {
			return list[i].b.CreatedAt.After(list[j].b.CreatedAt)
		})
		capN := tierCap(p, tier)
		for i, c := range list {
			b := c.b
			b.Tier = tier
			out.Tiers[b.BackupId] = tier
			if i < capN || b.BackupId == newestID {
				out.Keep = append(out.Keep, b)
			} else {
				out.Prune = append(out.Prune, b)
			}
		}
	}

	// Ensure newest is present even if somehow dropped.
	if !containsBackupID(out.Keep, newestID) {
		for i, b := range out.Prune {
			if b.BackupId == newestID {
				out.Keep = append(out.Keep, b)
				out.Prune = append(out.Prune[:i], out.Prune[i+1:]...)
				break
			}
		}
	}

	sortKeepPrune(&out)
	return out
}

func resolveAnchor(group []Backup, loc *time.Location, isAnchorDay func(time.Time) bool) string {
	var bestOnDay *Backup
	var newest *Backup
	for i := range group {
		b := &group[i]
		if newest == nil || b.CreatedAt.After(newest.CreatedAt) {
			newest = b
		}
		if isAnchorDay(b.CreatedAt) {
			if bestOnDay == nil || b.CreatedAt.After(bestOnDay.CreatedAt) {
				bestOnDay = b
			}
		}
	}
	if bestOnDay != nil {
		return bestOnDay.BackupId
	}
	if newest != nil {
		return newest.BackupId
	}
	return ""
}

func tierCap(p RetentionPolicy, tier BackupTier) int {
	// Config validation requires >= 1; treat <= 0 as unlimited so a zero-value
	// policy in tests cannot wipe every backup of a tier.
	n := 0
	switch tier {
	case TierWeekly:
		n = p.Weekly
	case TierMonthly:
		n = p.Monthly
	default:
		n = p.Daily
	}
	if n <= 0 {
		return math.MaxInt
	}
	return n
}

func containsBackupID(list []Backup, id string) bool {
	for _, b := range list {
		if b.BackupId == id {
			return true
		}
	}
	return false
}

func sortKeepPrune(out *RetentionPlan) {
	sort.Slice(out.Keep, func(i, j int) bool {
		return out.Keep[i].CreatedAt.After(out.Keep[j].CreatedAt)
	})
	sort.Slice(out.Prune, func(i, j int) bool {
		return out.Prune[i].CreatedAt.After(out.Prune[j].CreatedAt)
	})
	if out.Keep == nil {
		out.Keep = []Backup{}
	}
	if out.Prune == nil {
		out.Prune = []Backup{}
	}
}
