package domain

import "time"

// DailyBackupPoint is one day of aggregate backup size data.
type DailyBackupPoint struct {
	Date      string `json:"date"`
	SizeBytes int64  `json:"size_bytes"`
	Count     int    `json:"count"`
}

// BackupStats holds dashboard aggregates over backup metadata.
type BackupStats struct {
	TotalCount       int                `json:"total_count"`
	SuccessCount     int                `json:"success_count"`
	FailedCount      int                `json:"failed_count"`
	RunningCount     int                `json:"running_count"`
	TotalSizeBytes   int64              `json:"total_size_bytes"`
	LastSuccessAt    *time.Time         `json:"last_success_at,omitempty"`
	LastFailureAt    *time.Time         `json:"last_failure_at,omitempty"`
	LastSuccessID    string             `json:"last_success_id,omitempty"`
	LastFailureID    string             `json:"last_failure_id,omitempty"`
	DailySeries      []DailyBackupPoint `json:"daily_series"`
}

// RetentionPreview is the keep/prune split for the current GFS policy.
type RetentionPreview struct {
	Keep  []Backup `json:"keep"`
	Prune []Backup `json:"prune"`
}
