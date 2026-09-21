package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/backup"
	"github.com/sboy99/go-vault/internal/job"
	"github.com/sboy99/go-vault/pkg/logger"
)

// Scheduler runs cron-triggered backups and prune.
type Scheduler struct {
	cron    *cron.Cron
	cfg     *config.Config
	svc     *backup.Service
	runner  *job.Runner
	entryID cron.EntryID
}

func New(cfg *config.Config, svc *backup.Service, runner *job.Runner) (*Scheduler, error) {
	loc, err := time.LoadLocation(cfg.Schedule.Timezone)
	if err != nil {
		return nil, fmt.Errorf("load timezone %q: %w", cfg.Schedule.Timezone, err)
	}
	c := cron.New(cron.WithLocation(loc))
	return &Scheduler{cron: c, cfg: cfg, svc: svc, runner: runner}, nil
}

func (s *Scheduler) Start() error {
	id, err := s.cron.AddFunc(s.cfg.Schedule.Cron, func() {
		logger.Info("scheduled backup starting")
		_, err := s.runner.TryStart(job.TypeBackup, "", func(ctx context.Context) error {
			_, err := s.svc.BackupAndPrune(ctx)
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

func (s *Scheduler) Stop(ctx context.Context) {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-ctx.Done():
	}
}
