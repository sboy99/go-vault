import Link from "next/link";
import { notFound } from "next/navigation";
import { formatDistanceToNow } from "date-fns";
import { Download } from "lucide-react";
import { ApiError, getBackup, getConfig } from "@/lib/api/client";
import { formatBytes, formatDuration } from "@/lib/format";
import { fromPathId, toPathId } from "@/lib/route-id";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { BackupStatusBadge, BackupTierBadge } from "@/components/status-badge";
import { RestoreBackupButton } from "@/components/backups/restore-backup-button";
import { Card, CardTitle } from "@/components/ui/card";
import { CopyButton } from "@/components/ui/copy-button";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: pathId } = await params;
  const id = fromPathId(pathId);
  return { title: `Backup ${id.slice(0, 8)}` };
}

export default async function BackupDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: pathId } = await params;
  const id = fromPathId(pathId);

  let result;
  try {
    result = await getBackup(id);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }

  const configRes = await getConfig();
  const backup = result.data;
  const demo = result.demo || configRes.demo;
  const dbName = configRes.data.db.name;
  const canRestore = backup.status === "success";

  return (
    <div>
      <DemoBanner demo={demo} />
      <PageHeader
        title={backup.name}
        description={`Created ${formatDistanceToNow(new Date(backup.created_at), { addSuffix: true })}`}
        actions={
          <div className="flex flex-wrap items-center gap-2">
            <BackupStatusBadge status={backup.status} />
            {backup.tier ? <BackupTierBadge tier={backup.tier} /> : null}
            {canRestore ? (
              <>
                <a
                  href={`/api/backups/${toPathId(backup.id)}/download`}
                  aria-label="Download"
                  className="inline-flex size-8 items-center justify-center rounded-none border border-dashed border-edge bg-transparent text-foreground no-underline transition-colors duration-200 ease-out hover:text-brand"
                >
                  <Download className="size-4" strokeWidth={1.75} />
                </a>
                <RestoreBackupButton backupId={backup.id} dbName={dbName} />
              </>
            ) : null}
          </div>
        }
      />

      <div className="grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge sm:grid-cols-2 lg:grid-cols-4">
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Size
          </p>
          <p className="mt-2 font-mono text-[14px] text-foreground">
            {formatBytes(backup.size_bytes)}
          </p>
        </Card>
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Duration
          </p>
          <p className="mt-2 font-mono text-[14px] text-foreground">
            {formatDuration(backup.started_at, backup.finished_at)}
          </p>
        </Card>
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Format
          </p>
          <p className="mt-2 text-[14px] text-foreground">{backup.format}</p>
        </Card>
        <Card variant="flat">
          <p className="text-xs uppercase tracking-wide text-foreground-muted">
            Verified
          </p>
          <p className="mt-2 text-[14px] text-foreground">
            {backup.verified ? "Yes" : "No"}
          </p>
        </Card>
      </div>

      <Card
        variant="flat"
        className="mt-3 border border-dashed border-edge"
      >
        <div className="mb-3 flex items-center justify-between gap-2 border-b border-dashed border-edge pb-3">
          <CardTitle>Details</CardTitle>
          <CopyButton value={backup.id} label="Copy id" />
        </div>
        <dl className="space-y-2 text-[13px]">
          <Row label="Backup id" value={backup.id} mono />
          <Row label="Storage key" value={backup.storage_key} mono />
          <Row label="Database" value={backup.database_type} />
          <Row label="Storage" value={backup.storage_type} />
          {backup.tier ? (
            <div className="flex flex-col gap-0.5 sm:flex-row sm:gap-4">
              <dt className="w-28 shrink-0 text-foreground-muted">Type</dt>
              <dd className="text-foreground">
                <BackupTierBadge tier={backup.tier} />
              </dd>
            </div>
          ) : null}
          {backup.pg_version ? (
            <Row label="PG version" value={backup.pg_version} />
          ) : null}
          {backup.sha256 ? <Row label="SHA-256" value={backup.sha256} mono /> : null}
          <Row
            label="Started"
            value={
              backup.started_at
                ? new Date(backup.started_at).toLocaleString()
                : "—"
            }
          />
          <Row
            label="Finished"
            value={
              backup.finished_at
                ? new Date(backup.finished_at).toLocaleString()
                : "—"
            }
          />
        </dl>
      </Card>

      {backup.error ? (
        <Card
          variant="flat"
          className="mt-3 border border-dashed border-destructive/40 bg-destructive-muted"
        >
          <CardTitle className="text-destructive">Error</CardTitle>
          <p className="mt-1 font-mono text-[13px] text-destructive">
            {backup.error}
          </p>
        </Card>
      ) : null}

      <p className="mt-6">
        <Link
          href="/backups"
          className="text-[13px] text-brand no-underline hover:underline"
        >
          ← All backups
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
