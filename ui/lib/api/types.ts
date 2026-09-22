export type BackupStatus = "running" | "success" | "failed";
export type JobStatus = "queued" | "running" | "succeeded" | "failed";
export type JobType = "backup" | "restore";
export type BackupTier = "daily" | "weekly" | "monthly";

export type Backup = {
  id: string;
  name: string;
  storage_key: string;
  database_type: string;
  storage_type: string;
  status: BackupStatus;
  error?: string;
  started_at: string;
  finished_at?: string;
  created_at: string;
  size_bytes: number;
  sha256?: string;
  pg_version?: string;
  format: string;
  verified: boolean;
  tier?: BackupTier;
};

export type Job = {
  id: string;
  type: JobType;
  status: JobStatus;
  error?: string;
  backup_id?: string;
  created_at: string;
  started_at?: string;
  finished_at?: string;
};

export type DailyBackupPoint = {
  date: string;
  size_bytes: number;
  count: number;
};

export type BackupStats = {
  total_count: number;
  success_count: number;
  failed_count: number;
  running_count: number;
  total_size_bytes: number;
  last_success_at?: string;
  last_failure_at?: string;
  last_success_id?: string;
  last_failure_id?: string;
  daily_series: DailyBackupPoint[];
};

export type StatsResponse = {
  backups: BackupStats;
  job_running: boolean;
};

export type RedactedConfig = {
  app: { name: string; version: string };
  db: {
    type: string;
    host: string;
    port: number;
    name: string;
    sslmode: string;
  };
  storage: {
    type: string;
    dest?: string;
    cloud?: {
      type: string;
      region?: string;
      bucket?: string;
      endpoint?: string;
    };
  };
  schedule: { cron: string; timezone: string };
  retention: { daily: number; weekly: number; monthly: number };
  next_run?: string;
};

export type RetentionPreview = {
  keep: Backup[] | null;
  prune: Backup[] | null;
};

export type HealthStatus = {
  healthz: boolean;
  readyz: boolean;
  health_error?: string;
  ready_error?: string;
};

export type Envelope<T> = {
  success: boolean;
  data: T | null;
  error: string | null;
};

export type DataResult<T> = {
  data: T;
  demo: boolean;
};
