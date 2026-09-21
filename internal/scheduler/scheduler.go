package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/pkg/logger"
)

// Scheduler runs cron-triggered backups and prune.
type Scheduler struct {
	cron    *cron.Cron
	cfg     *config.Config
	backup  *app.BackupService
	jobs    *app.JobService
	entryID cron.EntryID
}

func New(cfg *config.Config, backup *app.BackupService, jobs *app.JobService) (*Scheduler, error) {
	loc, err := time.LoadLocation(cfg.Schedule.Timezone)
	if err != nil {
		return nil, fmt.Errorf("load timezone %q: %w", cfg.Schedule.Timezone, err)
	}
	c := cron.New(cron.WithLocation(loc))
	return &Scheduler{cron: c, cfg: cfg, backup: backup, jobs: jobs}, nil
}

func (s *Scheduler) Start() error {
	id, err := s.cron.AddFunc(s.cfg.Schedule.Cron, func() {
		logger.Info("scheduled backup starting")
		_, err := s.jobs.TryStart(domain.JobTypeBackup, "", func(ctx context.Context) error {
			_, err := s.backup.BackupAndPrune(ctx)
			return err
		})
		if err != nil {
			logger.Warn("scheduled backup skipped: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid schedule.cron %q: %w", s.cfg.Schedule.Cron, err)
	}
	s.entryID = id
	s.cron.Start()
	logger.Info("scheduler started cron=%s tz=%s", s.cfg.Schedule.Cron, s.cfg.Schedule.Timezone)
	return nil
}

// NextRun returns the next scheduled fire time, or zero if unknown.
func (s *Scheduler) NextRun() time.Time {
	if s == nil || s.cron == nil {
		return time.Time{}
	}
	entry := s.cron.Entry(s.entryID)
	return entry.Next
}

func (s *Scheduler) Stop(ctx context.Context) {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-ctx.Done():
	}
}
