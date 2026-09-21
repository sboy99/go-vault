import { Badge } from "@/components/ui/badge";
import type { BackupStatus, JobStatus } from "@/lib/api/types";

const backupLabels: Record<BackupStatus, string> = {
  success: "Success",
  failed: "Failed",
  running: "Running",
};

const jobLabels: Record<JobStatus, string> = {
  succeeded: "Succeeded",
  failed: "Failed",
  running: "Running",
  queued: "Queued",
};

export function BackupStatusBadge({ status }: { status: BackupStatus }) {
  const tone =
    status === "success"
      ? "success"
      : status === "failed"
        ? "error"
        : "running";
  return <Badge tone={tone}>{backupLabels[status]}</Badge>;
}

export function JobStatusBadge({ status }: { status: JobStatus }) {
  const tone =
    status === "succeeded"
      ? "success"
      : status === "failed"
        ? "error"
        : status === "running"
          ? "running"
          : "info";
  return <Badge tone={tone}>{jobLabels[status]}</Badge>;
}
