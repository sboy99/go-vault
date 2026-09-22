package domain

import "time"

// BackupTier is the rollup classification of a backup artifact.
type BackupTier string

const (
	TierDaily   BackupTier = "daily"
	TierWeekly  BackupTier = "weekly"
	TierMonthly BackupTier = "monthly"
)

// dateOnly returns midnight of t in loc.
func dateOnly(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

// weekSunday returns the Sunday that starts the Sun–Sat week containing t.
func weekSunday(t time.Time, loc *time.Location) time.Time {
	d := dateOnly(t, loc)
	return d.AddDate(0, 0, -int(d.Weekday()))
}

// weekCycleEnd returns the Saturday of the Sun–Sat week containing t.
func weekCycleEnd(t time.Time, loc *time.Location) time.Time {
	return weekSunday(t, loc).AddDate(0, 0, 6)
}

// monthCycleEnd returns the last calendar day of the month containing t.
func monthCycleEnd(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	firstNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
	return firstNext.AddDate(0, 0, -1)
}

// cycleClosed reports whether the cycle whose last day is endDay has ended
// relative to now (both interpreted in loc).
func cycleClosed(endDay, now time.Time, loc *time.Location) bool {
	return dateOnly(now, loc).After(dateOnly(endDay, loc))
}

func weekKey(t time.Time, loc *time.Location) string {
	return weekSunday(t, loc).Format("2006-01-02")
}

func monthKey(t time.Time, loc *time.Location) string {
	return t.In(loc).Format("2006-01")
}
