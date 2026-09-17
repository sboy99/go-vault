package metrics

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	BackupLastSuccessTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "govault_backup_last_success_timestamp",
		Help: "Unix timestamp of the last successful backup",
	})
	BackupDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "govault_backup_duration_seconds",
		Help:    "Duration of backup jobs in seconds",
		Buckets: prometheus.ExponentialBuckets(1, 2, 12),
	})
	BackupSizeBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "govault_backup_size_bytes",
		Help: "Size in bytes of the last successful backup",
	})
	BackupFailuresTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "govault_backup_failures_total",
		Help: "Total number of failed backups",
	})
	RestoreTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "govault_restore_total",
		Help: "Total number of restore attempts by status",
	}, []string{"status"})
	RetentionPrunedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "govault_retention_pruned_total",
		Help: "Total number of backups pruned by retention policy",
	})
)

var ready atomic.Bool

func SetReady(v bool) { ready.Store(v) }
func Ready() bool     { return ready.Load() }

func Handler() http.Handler {
	return promhttp.Handler()
}

func ObserveBackupSuccess(started time.Time, size int64) {
	BackupLastSuccessTimestamp.Set(float64(time.Now().Unix()))
	BackupDurationSeconds.Observe(time.Since(started).Seconds())
	BackupSizeBytes.Set(float64(size))
}

func ObserveBackupFailure() {
	BackupFailuresTotal.Inc()
}
