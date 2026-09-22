import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { listJobs } from "@/lib/api/client";
import { formatDuration } from "@/lib/format";
import { jobHref } from "@/lib/route-id";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { JobStatusBadge } from "@/components/status-badge";
import { JobsPoller } from "@/components/jobs/jobs-poller";
import { EmptyState } from "@/components/ui/empty-state";
import { Table, TBody, TD, TH, THead, TR } from "@/components/ui/table";

export const metadata = { title: "Jobs" };

export default async function JobsPage() {
  const { data, demo } = await listJobs({ limit: 50 });
  const running = data.some(
    (j) => j.status === "running" || j.status === "queued",
  );

  return (
    <div>
      <DemoBanner demo={demo} />
      <JobsPoller enabled={running} />
      <PageHeader
        title="Jobs"
        description="Async backup and restore operations. Only one job runs at a time."
      />

      {data.length === 0 ? (
        <EmptyState
          title="No jobs yet"
          description="Jobs appear when you create a backup or restore."
        />
      ) : (
        <Table>
          <THead>
            <TR>
              <TH>Job</TH>
              <TH>Type</TH>
              <TH>Status</TH>
              <TH>Duration</TH>
              <TH>Created</TH>
            </TR>
          </THead>
          <TBody>
            {data.map((job) => (
              <TR key={job.id}>
                <TD>
                  <Link
                    href={jobHref(job.id)}
                    className="font-mono text-[12px] text-foreground no-underline hover:text-brand"
                  >
                    {job.id}
                  </Link>
                  {job.backup_id ? (
                    <p className="max-w-[220px] truncate text-xs text-foreground-muted">
                      backup {job.backup_id.slice(-12)}
                    </p>
                  ) : null}
                </TD>
                <TD className="capitalize">{job.type}</TD>
                <TD>
                  <JobStatusBadge status={job.status} />
                </TD>
                <TD className="font-mono text-[12px]">
                  {formatDuration(job.started_at, job.finished_at)}
                </TD>
                <TD className="text-foreground-muted">
                  {formatDistanceToNow(new Date(job.created_at), {
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
