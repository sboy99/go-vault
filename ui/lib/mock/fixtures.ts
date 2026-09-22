import type {
  Backup,
  BackupStats,
  DailyBackupPoint,
  Job,
  RedactedConfig,
  RetentionPreview,
  StatsResponse,
} from "@/lib/api/types";

function daysAgo(n: number): string {
  const d = new Date();
  d.setUTCDate(d.getUTCDate() - n);
  return d.toISOString();
}

function dayKey(n: number): string {
  const d = new Date();
  d.setUTCDate(d.getUTCDate() - n);
  return d.toISOString().slice(0, 10);
}

const mockBackups: Backup[] = [
  {
    id: "2026-09-22T01:00:00.000000000Z_a1b2c3d4e5f6",
    name: "1726966800_app_backup.dump",
    storage_key: "1726966800_app_backup.dump",
    database_type: "POSTGRESQL",
    storage_type: "LOCAL",
    status: "success",
    started_at: daysAgo(0),
    finished_at: daysAgo(0),
    created_at: daysAgo(0),
    size_bytes: 48_291_840,
    sha256: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
    pg_version: "16.4",
    format: "custom",
    verified: true,
    tier: "daily",
  },
  {
    // Saturday of a closed week (Sun–Sat) in the open month.
    id: "2026-09-13T01:00:00.000000000Z_b2c3d4e5f6a7",
    name: "1726880400_app_backup.dump",
    storage_key: "1726880400_app_backup.dump",
    database_type: "POSTGRESQL",
    storage_type: "LOCAL",
    status: "success",
    started_at: "2026-09-13T01:00:00.000Z",
    finished_at: "2026-09-13T01:00:05.000Z",
    created_at: "2026-09-13T01:00:00.000Z",
    size_bytes: 47_185_920,
    sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    pg_version: "16.4",
    format: "custom",
    verified: true,
    tier: "weekly",
  },
  {
    id: "2026-09-20T01:00:00.000000000Z_c3d4e5f6a7b8",
    name: "1726794000_app_backup.dump",
    storage_key: "1726794000_app_backup.dump",
    database_type: "POSTGRESQL",
    storage_type: "LOCAL",
    status: "failed",
    error: "dump: connection refused",
    started_at: daysAgo(2),
    finished_at: daysAgo(2),
    created_at: daysAgo(2),
    size_bytes: 0,
    format: "custom",
    verified: false,
    tier: "daily",
  },
  {
    // Last calendar day of a closed month.
    id: "2026-08-31T01:00:00.000000000Z_d4e5f6a7b8c9",
    name: "1726707600_app_backup.dump",
    storage_key: "1726707600_app_backup.dump",
    database_type: "POSTGRESQL",
    storage_type: "LOCAL",
    status: "success",
    started_at: "2026-08-31T01:00:00.000Z",
    finished_at: "2026-08-31T01:00:05.000Z",
    created_at: "2026-08-31T01:00:00.000Z",
    size_bytes: 46_137_344,
    sha256: "2c624232cdd221771294dfbb310aca000a0df6ac8b66b696d90ef06fdefb64a3",
    pg_version: "16.4",
    format: "custom",
    verified: true,
    tier: "monthly",
  },
];

const mockJobs: Job[] = [
  {
    id: "job-demo-1",
    type: "backup",
    status: "succeeded",
    created_at: daysAgo(0),
    started_at: daysAgo(0),
    finished_at: daysAgo(0),
  },
  {
    id: "job-demo-2",
    type: "restore",
    status: "failed",
    backup_id: mockBackups[1].id,
    error: "confirm must match the configured database name",
    created_at: daysAgo(1),
    started_at: daysAgo(1),
    finished_at: daysAgo(1),
  },
  {
    id: "job-demo-3",
    type: "backup",
    status: "succeeded",
    created_at: daysAgo(1),
    started_at: daysAgo(1),
    finished_at: daysAgo(1),
  },
];

function buildDailySeries(): DailyBackupPoint[] {
  return Array.from({ length: 30 }, (_, i) => {
    const n = 29 - i;
    const match = mockBackups.find((b) => b.created_at.slice(0, 10) === dayKey(n));
    return {
      date: dayKey(n),
      size_bytes: match?.status === "success" ? match.size_bytes : 0,
      count: match?.status === "success" ? 1 : 0,
    };
  });
}

const mockStats: BackupStats = {
  total_count: mockBackups.length,
  success_count: mockBackups.filter((b) => b.status === "success").length,
  failed_count: mockBackups.filter((b) => b.status === "failed").length,
  running_count: 0,
  total_size_bytes: mockBackups
    .filter((b) => b.status === "success")
    .reduce((sum, b) => sum + b.size_bytes, 0),
  last_success_at: mockBackups[0].created_at,
  last_success_id: mockBackups[0].id,
  last_failure_at: mockBackups[2].created_at,
  last_failure_id: mockBackups[2].id,
  daily_series: buildDailySeries(),
};

const mockConfig: RedactedConfig = {
  app: { name: "go-vault", version: "0.2.0" },
  db: {
    type: "POSTGRESQL",
    host: "db.internal",
    port: 5432,
    name: "app",
    sslmode: "require",
  },
  storage: {
    type: "LOCAL",
    dest: "/var/lib/go-vault/backups",
  },
  schedule: { cron: "0 2 * * *", timezone: "UTC" },
  retention: { daily: 7, weekly: 4, monthly: 12 },
  next_run: new Date(Date.now() + 6 * 60 * 60 * 1000).toISOString(),
};

export const mocks = {
  backups: mockBackups,
  jobs: mockJobs,
  stats: { backups: mockStats, job_running: false } satisfies StatsResponse,
  config: mockConfig,
  retention: {
    keep: mockBackups.filter((b) => b.status === "success"),
    prune: [],
  } satisfies RetentionPreview,
  health: { healthz: true, readyz: true },
};

export function mockBackupById(id: string): Backup | undefined {
  return mockBackups.find((b) => b.id === id);
}

export function mockJobById(id: string): Job | undefined {
  return mockJobs.find((j) => j.id === id);
}
