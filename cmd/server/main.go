package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sboy99/go-vault/internal/bootstrap"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/internal/scheduler"
	transporthttp "github.com/sboy99/go-vault/internal/transport/http"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	logger.Init(logger.INFO)

	if err := config.Load(); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}

	cfg := config.GetConfig()
	c, err := bootstrap.New(cfg)
	if err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := c.Close(); err != nil {
			logger.Error("%s", err.Error())
		}
	}()

	reconcileOnStart(c)

	sched, err := scheduler.New(cfg, c.Backup, c.Jobs)
	if err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
	if err := sched.Start(); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}

	metrics.SetReady(true)
	handler := transporthttp.NewRouter(transporthttp.RouterDeps{
		Backup:   c.Backup,
		Jobs:     c.Jobs,
		APIToken: cfg.API.Token,
	})
	server := &http.Server{
		Addr:              cfg.API.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api listening on %s", cfg.API.Addr)
		errCh <- server.ListenAndServe()
	}()

	waitForShutdown(errCh)
	if err := shutdown(cfg, sched, c, server); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
}

func reconcileOnStart(c *bootstrap.Container) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Backup.Reconcile(ctx); err != nil {
		logger.Warn("reconcile: %v", err)
	}
}

func waitForShutdown(errCh <-chan error) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("received signal %s, shutting down", sig.String())
	case err := <-errCh:
		if err != nil && err.Error() != "http: Server closed" {
			logger.Error("api error: %v", err)
		}
	}
}

func shutdown(cfg *config.Config, sched *scheduler.Scheduler, c *bootstrap.Container, server *http.Server) error {
	metrics.SetReady(false)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Runtime.ShutdownTimeout)
	defer cancel()
	sched.Stop(ctx)
	c.Jobs.Cancel()
	return server.Shutdown(ctx)
}
