import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import {
  getConfig,
  getStats,
  listBackups,
  listJobs,
} from "@/lib/api/client";
import type {
  Backup,
  Job,
  RedactedConfig,
  StatsResponse,
} from "@/lib/api/types";
import { formatBytes } from "@/lib/format";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { Card, CardTitle } from "@/components/ui/card";
import { SizeChart } from "@/components/overview/size-chart";
import { RunningJobBanner } from "@/components/overview/running-job-banner";
import { BackupStatusBadge, JobStatusBadge } from "@/components/status-badge";
import { EmptyState } from "@/components/ui/empty-state";

export const metadata = { title: "Overview" };

async function safe<T>(
  fn: () => Promise<{ data: T; demo: boolean }>,
  fallback: T,
): Promise<{ data: T; demo: boolean; error?: string }> {
  try {
    return await fn();
  } catch (err) {
    return {
      data: fallback,
      demo: true,
      error: err instanceof Error ? err.message : "failed to load",
    };
  }
}

const emptyStats: StatsResponse = {
  backups: {
    total_count: 0,
    success_count: 0,
    failed_count: 0,
    running_count: 0,
    total_size_bytes: 0,
    daily_series: [],
  },
  job_running: false,
};

const emptyConfig: RedactedConfig = {
  app: { name: "go-vault", version: "" },
  db: { type: "", host: "", port: 0, name: "", sslmode: "" },
  storage: { type: "" },
  schedule: { cron: "", timezone: "" },
  retention: { daily: 0, weekly: 0, monthly: 0 },
};

export default async function OverviewPage() {
  const [statsRes, configRes, backupsRes, jobsRes] = await Promise.all([
    safe(getStats, emptyStats),
    safe(getConfig, emptyConfig),
    safe(() => listBackups({ limit: 5 }), [] as Backup[]),
    safe(() => listJobs({ limit: 5 }), [] as Job[]),
  ]);

  const demo =
    statsRes.demo || configRes.demo || backupsRes.demo || jobsRes.demo;
  const stats = statsRes.data.backups;
  const nextRun = configRes.data.next_run;

  return (
    <div>
      <DemoBanner demo={demo} />
      <PageHeader
        title="Overview"
        description="Live picture of PostgreSQL backups, schedule health, and recent jobs."
      />
      <RunningJobBanner running={statsRes.data.job_running} />

      <div className="grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Last success"
          value={
            stats.last_success_at
              ? formatDistanceToNow(new Date(stats.last_success_at), {
                  addSuffix: true,
                })
              : "Never"
          }
          hint={
            stats.last_success_id
              ? `id ${stats.last_success_id.slice(-8)}`
              : undefined
          }
        />
        <StatCard
          label="Artifacts"
          value={String(stats.total_count)}
          hint={`${stats.success_count} ok · ${stats.failed_count} failed`}
        />
        <StatCard
          label="Total size"
          value={formatBytes(stats.total_size_bytes)}
          hint={`${stats.running_count} running`}
        />
        <StatCard
          label="Next run"
          value={
            nextRun
              ? formatDistanceToNow(new Date(nextRun), { addSuffix: true })
              : "—"
          }
          hint={configRes.data.schedule.cron || undefined}
        />
      </div>

      <div className="mt-3 grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge lg:grid-cols-5">
        <Card variant="flat" className="lg:col-span-3">
          <CardTitle>Backup size · 30 days</CardTitle>
          <p className="mt-1 mb-3 text-sm text-foreground-muted">
            Successful dump sizes day by day.
          </p>
          {stats.daily_series.length > 0 ? (
            <SizeChart series={stats.daily_series} />
          ) : (
            <EmptyState title="No series data" className="py-8" />
          )}
        </Card>

        <Card variant="flat" className="lg:col-span-2">
          <CardTitle>Recent activity</CardTitle>
          <p className="mt-1 mb-3 text-sm text-foreground-muted">
            Latest jobs and backups.
          </p>
          {jobsRes.data.length === 0 ? (
            <EmptyState
              title="No jobs yet"
              description="Trigger a backup to see activity."
            />
          ) : (
            <ul className="divide-y divide-dashed divide-edge">
              {jobsRes.data.map((job) => (
                <li
                  key={job.id}
                  className="flex items-center justify-between gap-2 py-1.5"
                >
                  <div className="min-w-0">
                    <Link
                      href={`/jobs/${job.id}`}
                      className="block truncate text-[13px] text-foreground no-underline hover:text-brand"
                    >
                      {job.type} · {job.id.slice(0, 8)}
                    </Link>
                    <p className="text-xs text-foreground-muted">
                      {formatDistanceToNow(new Date(job.created_at), {
                        addSuffix: true,
                      })}
                    </p>
                  </div>
                  <JobStatusBadge status={job.status} />
                </li>
              ))}
            </ul>
          )}
        </Card>
      </div>

      <Card variant="flat" className="mt-3 border border-dashed border-edge">
        <div className="mb-3 flex items-center justify-between border-b border-dashed border-edge pb-3">
          <CardTitle>Recent backups</CardTitle>
          <Link
            href="/backups"
            className="text-[13px] text-brand no-underline hover:underline"
          >
            View all
          </Link>
        </div>
        {backupsRes.data.length === 0 ? (
          <EmptyState
            title="No backups"
            description="Create your first dump from Backups."
          />
        ) : (
          <ul className="divide-y divide-dashed divide-edge">
            {backupsRes.data.map((b) => (
              <li
                key={b.id}
                className="flex flex-wrap items-center justify-between gap-2 py-1.5"
              >
                <div className="min-w-0">
                  <Link
                    href={`/backups/${encodeURIComponent(b.id)}`}
                    className="truncate font-mono text-[13px] text-foreground no-underline hover:text-brand"
                  >
                    {b.name}
                  </Link>
                  <p className="text-xs text-foreground-muted">
                    {formatBytes(b.size_bytes)} ·{" "}
                    {formatDistanceToNow(new Date(b.created_at), {
                      addSuffix: true,
                    })}
                  </p>
                </div>
                <BackupStatusBadge status={b.status} />
              </li>
            ))}
          </ul>
        )}
      </Card>
    </div>
  );
}

function StatCard({
  label,
  value,
  hint,
}: {
  label: string;
  value: string;
  hint?: string;
}) {
  return (
    <Card variant="flat">
      <p className="text-xs uppercase tracking-wide text-foreground-muted">
        {label}
      </p>
      <p className="mt-2 font-heading text-xl font-semibold text-foreground">
        {value}
      </p>
      {hint ? (
        <p className="mt-1 truncate text-xs text-foreground-muted">{hint}</p>
      ) : null}
    </Card>
  );
}
