import Link from "next/link";
import { notFound } from "next/navigation";
import { formatDistanceToNow } from "date-fns";
import { ApiError, getJob } from "@/lib/api/client";
import { formatDuration } from "@/lib/format";
import { backupHref, fromPathId } from "@/lib/route-id";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { JobStatusBadge } from "@/components/status-badge";
import { JobsPoller } from "@/components/jobs/jobs-poller";
import { Card, CardTitle } from "@/components/ui/card";
import { CopyButton } from "@/components/ui/copy-button";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: pathId } = await params;
  const id = fromPathId(pathId);
  return { title: `Job ${id.slice(0, 8)}` };
}

function titleCase(s: string) {
  if (!s) return s;
  return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
}

export default async function JobDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: pathId } = await params;
  const id = fromPathId(pathId);
  let result;
  try {
    result = await getJob(id);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }

  const job = result.data;
  const active = job.status === "running" || job.status === "queued";

  return (
    <div>
      <DemoBanner demo={result.demo} />
      <JobsPoller enabled={active} />
      <PageHeader
        title={`${titleCase(job.type)} job`}
        description={`Created ${formatDistanceToNow(new Date(job.created_at), { addSuffix: true })}`}
        actions={<JobStatusBadge status={job.status} />}
      />

      <div className="grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge sm:grid-cols-2 lg:grid-cols-3">
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Duration
          </p>
          <p className="mt-2 font-mono text-[14px] text-foreground">
            {formatDuration(job.started_at, job.finished_at)}
          </p>
        </Card>
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Started
          </p>
          <p className="mt-2 text-[14px] text-foreground">
            {job.started_at
              ? new Date(job.started_at).toLocaleString()
              : "—"}
          </p>
        </Card>
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Finished
          </p>
          <p className="mt-2 text-[14px] text-foreground">
            {job.finished_at
              ? new Date(job.finished_at).toLocaleString()
              : "—"}
          </p>
        </Card>
      </div>

      <Card
        variant="flat"
        className="mt-3 border border-dashed border-edge"
      >
        <div className="mb-3 flex items-center justify-between gap-2 border-b border-dashed border-edge pb-3">
          <CardTitle>Details</CardTitle>
          <CopyButton value={job.id} label="Copy id" />
        </div>
        <dl className="space-y-2 text-[13px]">
          <Row label="Job id" value={job.id} mono />
          <Row label="Type" value={titleCase(job.type)} />
          {job.backup_id ? (
            <div className="flex flex-col gap-0.5 sm:flex-row sm:gap-4">
              <dt className="w-28 shrink-0 text-foreground-muted">Backup</dt>
              <dd>
                <Link
                  href={backupHref(job.backup_id)}
                  className="break-all font-mono text-[12px] text-brand no-underline hover:underline"
                >
                  {job.backup_id}
                </Link>
              </dd>
            </div>
          ) : null}
        </dl>
      </Card>

      {job.error ? (
        <Card
          variant="flat"
          className="mt-3 border border-dashed border-destructive/40 bg-destructive-muted"
        >
          <CardTitle className="text-destructive">Error</CardTitle>
          <p className="mt-1 font-mono text-[13px] text-destructive">
            {job.error}
          </p>
        </Card>
      ) : null}

      <p className="mt-6">
        <Link
          href="/jobs"
          className="text-[13px] text-brand no-underline hover:underline"
        >
          ← All jobs
        </Link>
      </p>
    </div>
  );
}

function Row({
  label,
  value,
  mono,
}: {
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex flex-col gap-0.5 sm:flex-row sm:gap-4">
      <dt className="w-28 shrink-0 text-foreground-muted">{label}</dt>
      <dd
        className={`break-all text-foreground ${mono ? "font-mono text-[12px]" : ""}`}
      >
        {value}
      </dd>
    </div>
  );
}
