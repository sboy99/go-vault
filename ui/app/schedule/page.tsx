import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { getConfig, getRetention } from "@/lib/api/client";
import { describeCron, formatBytes } from "@/lib/format";
import { backupHref } from "@/lib/route-id";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { Card, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { EmptyState } from "@/components/ui/empty-state";
import { Table, TBody, TD, TH, THead, TR } from "@/components/ui/table";

export const metadata = { title: "Schedule" };

export default async function SchedulePage() {
  const [configRes, retentionRes] = await Promise.all([
    getConfig(),
    getRetention(),
  ]);
  const demo = configRes.demo || retentionRes.demo;
  const cfg = configRes.data;
  const keep = retentionRes.data.keep ?? [];
  const prune = retentionRes.data.prune ?? [];

  return (
    <div>
      <DemoBanner demo={demo} />
      <PageHeader
        title="Schedule"
        description="Cron-triggered backups and GFS retention policy."
      />

      <div className="relative grid gap-0 border border-dashed border-edge lg:grid-cols-[2fr_1fr]">
        <Card variant="flat">
          <CardTitle className="mb-3 border-b border-dashed border-edge pb-3">
            Cron
          </CardTitle>
          <p className="font-mono text-[18px] text-foreground">
            {cfg.schedule.cron}
          </p>
          <p className="mt-1 text-sm text-foreground-muted">
            {describeCron(cfg.schedule.cron)}
          </p>
          <dl className="mt-3 grid gap-2 border-t border-dashed border-edge pt-3 text-[13px] sm:grid-cols-2">
            <div>
              <dt className="text-xs uppercase tracking-wide text-foreground-muted">
                Timezone
              </dt>
              <dd className="text-foreground">{cfg.schedule.timezone}</dd>
            </div>
            <div>
              <dt className="text-xs uppercase tracking-wide text-foreground-muted">
                Next run
              </dt>
              <dd className="text-foreground">
                {cfg.next_run
                  ? `${formatDistanceToNow(new Date(cfg.next_run), { addSuffix: true })} · ${new Date(cfg.next_run).toLocaleString()}`
                  : "—"}
              </dd>
            </div>
          </dl>
        </Card>

        <div
          aria-hidden
          className="pointer-events-none absolute inset-y-0 left-[66.666%] z-10 hidden border-l border-dashed border-edge lg:block"
        />

        <Card
          variant="flat"
          className="border-dashed border-edge max-lg:border-t"
        >
          <CardTitle className="mb-3 border-b border-dashed border-edge pb-3">
            GFS retention
          </CardTitle>
          <ul className="space-y-2 text-[13px]">
            <li className="flex justify-between border-b border-dashed border-edge pb-2">
              <span className="text-foreground-muted">Daily</span>
              <span className="text-foreground">{cfg.retention.daily}</span>
            </li>
            <li className="flex justify-between border-b border-dashed border-edge pb-2">
              <span className="text-foreground-muted">Weekly</span>
              <span className="text-foreground">{cfg.retention.weekly}</span>
            </li>
            <li className="flex justify-between">
              <span className="text-foreground-muted">Monthly</span>
              <span className="text-foreground">{cfg.retention.monthly}</span>
            </li>
          </ul>
        </Card>
      </div>

      <div className="mt-3 grid gap-0 divide-y divide-dashed divide-edge border border-dashed border-edge">
        <RetentionTable
          title="Keep"
          tone="success"
          empty="Nothing to keep yet"
          rows={keep}
        />
        <RetentionTable
          title="Would prune"
          tone="warning"
          empty="Nothing queued for prune"
          rows={prune}
        />
      </div>
    </div>
  );
}

function RetentionTable({
  title,
  tone,
  empty,
  rows,
}: {
  title: string;
  tone: "success" | "warning";
  empty: string;
  rows: Array<{
    id: string;
    name: string;
    size_bytes: number;
    created_at: string;
  }>;
}) {
  return (
    <Card variant="flat">
      <div className="mb-3 flex items-center gap-2 border-b border-dashed border-edge pb-3">
        <CardTitle>{title}</CardTitle>
        <Badge tone={tone}>{rows.length}</Badge>
      </div>
      {rows.length === 0 ? (
        <EmptyState title={empty} className="py-8" />
      ) : (
        <Table className="border-0">
          <THead>
            <TR>
              <TH>Name</TH>
              <TH>Size</TH>
              <TH>Created</TH>
            </TR>
          </THead>
          <TBody>
            {rows.map((b) => (
              <TR key={b.id}>
                <TD>
                  <Link
                    href={backupHref(b.id)}
                    className="font-mono text-[12px] text-brand no-underline hover:underline"
                  >
                    {b.name}
                  </Link>
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
    </Card>
  );
}
