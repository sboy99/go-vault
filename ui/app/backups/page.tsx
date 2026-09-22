import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { listBackups } from "@/lib/api/client";
import { formatBytes } from "@/lib/format";
import { backupHref } from "@/lib/route-id";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { BackupStatusBadge, BackupTierBadge } from "@/components/status-badge";
import { RunBackupButton } from "@/components/backups/run-backup-button";
import { EmptyState } from "@/components/ui/empty-state";
import { Table, TBody, TD, TH, THead, TR } from "@/components/ui/table";

export const metadata = { title: "Backups" };

export default async function BackupsPage() {
  const { data, demo } = await listBackups({ limit: 100 });

  return (
    <div>
      <DemoBanner demo={demo} />
      <PageHeader
        title="Backups"
        description="PostgreSQL dump artifacts on local disk or object storage."
        actions={<RunBackupButton />}
      />

      {data.length === 0 ? (
        <EmptyState
          title="No backups yet"
          description="Run a backup to create the first dump."
          action={<RunBackupButton />}
        />
      ) : (
        <Table>
          <THead>
            <TR>
              <TH>Name</TH>
              <TH>Status</TH>
              <TH>Type</TH>
              <TH>Size</TH>
              <TH>Created</TH>
            </TR>
          </THead>
          <TBody>
            {data.map((b) => (
              <TR key={b.id}>
                <TD>
                  <Link
                    href={backupHref(b.id)}
                    className="font-mono text-[12px] text-foreground no-underline hover:text-brand"
                  >
                    {b.name}
                  </Link>
                  <p className="max-w-[280px] truncate text-xs text-foreground-muted">
                    {b.id}
                  </p>
                </TD>
                <TD>
                  <BackupStatusBadge status={b.status} />
                </TD>
                <TD>
                  {b.tier ? <BackupTierBadge tier={b.tier} /> : "—"}
                </TD>
                <TD className="font-mono text-[12px]">
                  {formatBytes(b.size_bytes)}
                </TD>
                <TD className="text-foreground-muted">
                  {formatDistanceToNow(new Date(b.created_at), {
                    addSuffix: true,
                  })}
                </TD>
              </TR>
            ))}
          </TBody>
        </Table>
      )}
    </div>
  );
}
