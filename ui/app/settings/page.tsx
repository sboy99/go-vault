import { getConfig, getHealth } from "@/lib/api/client";
import { PageHeader } from "@/components/layout/page-header";
import { DemoBanner } from "@/components/demo-banner";
import { Card, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { CopyButton } from "@/components/ui/copy-button";

export const metadata = { title: "Settings" };

export default async function SettingsPage() {
  const [configRes, healthRes] = await Promise.all([getConfig(), getHealth()]);
  const demo = configRes.demo || healthRes.demo;
  const cfg = configRes.data;
  const health = healthRes.data;
  const metricsUrl =
    process.env.GOVAULT_METRICS_PUBLIC_URL ||
    `${(process.env.GOVAULT_API_URL || "http://127.0.0.1:8080").replace(/\/$/, "")}/metrics`;
  const metricsHostSafe = !metricsUrl.includes("@");

  return (
    <div>
      <DemoBanner demo={demo} />
      <PageHeader
        title="Settings"
        description="Redacted connection and storage configuration. Secrets never leave the Go server."
      />

      <div className="grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge sm:grid-cols-2">
        <Card variant="flat">
          <div className="flex items-center justify-between">
            <CardTitle>Health</CardTitle>
            <Badge tone={health.healthz ? "success" : "error"}>
              {health.healthz ? "Ok" : "Down"}
            </Badge>
          </div>
          <p className="mt-2 text-sm text-foreground-muted">
            {health.health_error || "/healthz responding"}
          </p>
        </Card>
        <Card variant="flat">
          <div className="flex items-center justify-between">
            <CardTitle>Readiness</CardTitle>
            <Badge tone={health.readyz ? "success" : "warning"}>
              {health.readyz ? "Ready" : "Not ready"}
            </Badge>
          </div>
          <p className="mt-2 text-sm text-foreground-muted">
            {health.ready_error || "Database reachable"}
          </p>
        </Card>
      </div>

      <div className="mt-3 grid gap-0 divide-x divide-y divide-dashed divide-edge border border-dashed border-edge lg:grid-cols-2">
        <Card variant="flat">
          <CardTitle>Application</CardTitle>
          <dl className="mt-3 space-y-2 text-[13px]">
            <Row label="Name" value={cfg.app.name} />
            <Row label="Version" value={cfg.app.version} />
          </dl>
        </Card>

        <Card variant="flat">
          <CardTitle>Database</CardTitle>
          <dl className="mt-3 space-y-2 text-[13px]">
            <Row label="Type" value={cfg.db.type} />
            <Row label="Host" value={cfg.db.host} />
            <Row label="Port" value={String(cfg.db.port)} />
            <Row label="Name" value={cfg.db.name} />
            <Row label="SSL mode" value={cfg.db.sslmode} />
          </dl>
        </Card>

        <Card variant="flat">
          <CardTitle>Storage</CardTitle>
          <dl className="mt-3 space-y-2 text-[13px]">
            <Row label="Type" value={cfg.storage.type} />
            {cfg.storage.dest ? (
              <Row label="Dest" value={cfg.storage.dest} mono />
            ) : null}
            {cfg.storage.cloud ? (
              <>
                <Row label="Cloud" value={cfg.storage.cloud.type} />
                {cfg.storage.cloud.region ? (
                  <Row label="Region" value={cfg.storage.cloud.region} />
                ) : null}
                {cfg.storage.cloud.bucket ? (
                  <Row label="Bucket" value={cfg.storage.cloud.bucket} />
                ) : null}
                {cfg.storage.cloud.endpoint ? (
                  <Row
                    label="Endpoint"
                    value={cfg.storage.cloud.endpoint}
                    mono
                  />
                ) : null}
              </>
            ) : null}
          </dl>
        </Card>

        <Card variant="flat">
          <div className="flex items-center justify-between">
            <CardTitle>Metrics</CardTitle>
            {metricsHostSafe ? (
              <CopyButton value={metricsUrl} label="Copy URL" />
            ) : null}
          </div>
          <p className="mt-2 text-sm text-foreground-muted">
            Prometheus exposition endpoint on the go-vault API. Requires the same
            bearer token as other authenticated routes. Set{" "}
            <code className="font-mono text-[12px]">GOVAULT_METRICS_PUBLIC_URL</code>{" "}
            if the API is only reachable via an internal hostname.
          </p>
          {metricsHostSafe ? (
            <a
              href={metricsUrl}
              target="_blank"
              rel="noreferrer"
              className="mt-3 inline-block break-all font-mono text-[12px] text-brand no-underline hover:underline"
            >
              {metricsUrl}
            </a>
          ) : (
            <p className="mt-3 text-sm text-destructive">
              Metrics URL omitted — API URL contains credentials.
            </p>
          )}
        </Card>
      </div>
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
