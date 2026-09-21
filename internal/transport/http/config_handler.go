package http

import (
	"time"

	"github.com/sboy99/go-vault/internal/config"
	nethttp "net/http"
)

type redactedConfig struct {
	App       redactedApp       `json:"app"`
	DB        redactedDB        `json:"db"`
	Storage   redactedStorage   `json:"storage"`
	Schedule  redactedSchedule  `json:"schedule"`
	Retention redactedRetention `json:"retention"`
	NextRun   *time.Time        `json:"next_run,omitempty"`
}

type redactedApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type redactedDB struct {
	Type    string `json:"type"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Name    string `json:"name"`
	SSLMode string `json:"sslmode"`
}

type redactedStorage struct {
	Type  string              `json:"type"`
	Dest  string              `json:"dest,omitempty"`
	Cloud *redactedCloudStore `json:"cloud,omitempty"`
}

type redactedCloudStore struct {
	Type     string `json:"type"`
	Region   string `json:"region,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type redactedSchedule struct {
	Cron     string `json:"cron"`
	Timezone string `json:"timezone"`
}

type redactedRetention struct {
	Daily   int `json:"daily"`
	Weekly  int `json:"weekly"`
	Monthly int `json:"monthly"`
}

func (r *router) handleGetConfig(w nethttp.ResponseWriter, req *nethttp.Request) {
	if r.cfg == nil {
		writeError(w, nethttp.StatusInternalServerError, "config unavailable")
		return
	}
	writeOK(w, projectConfig(r.cfg, r.nextRun))
}

func projectConfig(cfg *config.Config, nextRunFn func() time.Time) redactedConfig {
	out := redactedConfig{
		App: redactedApp{
			Name:    cfg.App.Name,
			Version: cfg.App.Version,
		},
		DB: redactedDB{
			Type:    string(cfg.DB.Type),
			Host:    cfg.DB.Host,
			Port:    cfg.DB.Port,
			Name:    cfg.DB.Name,
			SSLMode: cfg.DB.SSLMode,
		},
		Storage: redactedStorage{
			Type: string(cfg.Storage.Type),
			Dest: cfg.Storage.Dest,
		},
		Schedule: redactedSchedule{
			Cron:     cfg.Schedule.Cron,
			Timezone: cfg.Schedule.Timezone,
		},
		Retention: redactedRetention{
			Daily:   cfg.Retention.Daily,
			Weekly:  cfg.Retention.Weekly,
			Monthly: cfg.Retention.Monthly,
		},
	}
	if cfg.Storage.Type == config.CLOUD {
		out.Storage.Cloud = &redactedCloudStore{
			Type:     string(cfg.Storage.Cloud.Type),
			Region:   cfg.Storage.Cloud.AWS.Region,
			Bucket:   cfg.Storage.Cloud.AWS.BucketName,
			Endpoint: cfg.Storage.Cloud.AWS.Endpoint,
		}
	}
	if nextRunFn != nil {
		t := nextRunFn()
		if !t.IsZero() {
			out.NextRun = &t
		}
	}
	return out
}
