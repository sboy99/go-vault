package domain

import (
	"fmt"
	"sort"
	"time"
)

// RetentionPolicy is a GFS retention policy.
type RetentionPolicy struct {
	Daily   int
	Weekly  int
	Monthly int
}

// Select returns backups to keep and prune according to GFS rules.
// Only successful backups are considered. The newest successful backup is always kept.
func Select(backups []Backup, p RetentionPolicy, now time.Time, loc *time.Location) (keep, prune []Backup) {
	if loc == nil {
		loc = time.UTC
	}
	now = now.In(loc)

	var success []Backup
	for _, b := range backups {
		if b.Status != StatusSuccess {
			continue
		}
		success = append(success, b)
	}
	if len(success) == 0 {
		return []Backup{}, []Backup{}
	}

	sort.Slice(success, func(i, j int) bool {
		return success[i].CreatedAt.After(success[j].CreatedAt)
	})

	keepSet := map[string]Backup{}
	keepSet[success[0].BackupId] = success[0]

	dailySlots := map[string]Backup{}
	weeklySlots := map[string]Backup{}
	monthlySlots := map[string]Backup{}

	for _, b := range success {
		t := b.CreatedAt.In(loc)
		dayKey := t.Format("2006-01-02")
		year, week := t.ISOWeek()
		weekKey := fmt.Sprintf("%04d-W%02d", year, week)
		monthKey := t.Format("2006-01")

		if existing, ok := dailySlots[dayKey]; !ok || b.CreatedAt.After(existing.CreatedAt) {
			dailySlots[dayKey] = b
		}
		if existing, ok := weeklySlots[weekKey]; !ok || b.CreatedAt.After(existing.CreatedAt) {
			weeklySlots[weekKey] = b
		}
		if existing, ok := monthlySlots[monthKey]; !ok || b.CreatedAt.After(existing.CreatedAt) {
			monthlySlots[monthKey] = b
		}
	}

	for key, b := range dailySlots {
		if inLastNDays(key, now, loc, p.Daily) {
			keepSet[b.BackupId] = b
		}
	}
	for key, b := range weeklySlots {
		if inLastNWeeks(key, now, loc, p.Weekly) {
			keepSet[b.BackupId] = b
		}
	}
	for key, b := range monthlySlots {
		if inLastNMonths(key, now, loc, p.Monthly) {
			keepSet[b.BackupId] = b
		}
	}

	for _, b := range success {
		if _, ok := keepSet[b.BackupId]; ok {
			keep = append(keep, b)
		} else {
			prune = append(prune, b)
		}
	}

	if len(keep) == 0 && len(success) > 0 {
		keep = []Backup{success[0]}
		prune = append([]Backup(nil), success[1:]...)
	}

	sort.Slice(keep, func(i, j int) bool { return keep[i].CreatedAt.After(keep[j].CreatedAt) })
	sort.Slice(prune, func(i, j int) bool { return prune[i].CreatedAt.After(prune[j].CreatedAt) })
	if keep == nil {
		keep = []Backup{}
	}
	if prune == nil {
		prune = []Backup{}
	}
	return keep, prune
}

func inLastNDays(dayKey string, now time.Time, loc *time.Location, n int) bool {
	if n <= 0 {
		return false
	}
	t, err := time.ParseInLocation("2006-01-02", dayKey, loc)
	if err != nil {
		return false
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -(n - 1))
	return !t.Before(start)
}

func inLastNWeeks(weekKey string, now time.Time, loc *time.Location, n int) bool {
	if n <= 0 {
		return false
	}
	allowed := map[string]struct{}{}
	t := now.In(loc)
	for i := 0; i < n; i++ {
		d := t.AddDate(0, 0, -7*i)
		y, w := d.ISOWeek()
		allowed[fmt.Sprintf("%04d-W%02d", y, w)] = struct{}{}
	}
	_, ok := allowed[weekKey]
	return ok
}

func inLastNMonths(monthKey string, now time.Time, loc *time.Location, n int) bool {
	if n <= 0 {
		return false
	}
	allowed := map[string]struct{}{}
	year, month, _ := now.In(loc).Date()
	for i := 0; i < n; i++ {
		m := time.Date(year, month, 1, 0, 0, 0, 0, loc).AddDate(0, -i, 0)
		allowed[m.Format("2006-01")] = struct{}{}
	}
	_, ok := allowed[monthKey]
	return ok
}
