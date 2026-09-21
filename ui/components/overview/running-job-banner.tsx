"use client";

import { useJobPoll } from "@/hooks/use-job-poll";
import { Badge } from "@/components/ui/badge";

export function RunningJobBanner({ running }: { running: boolean }) {
  useJobPoll(running);
  if (!running) return null;

  return (
    <div className="mb-4 flex items-center gap-2 rounded-md border border-info/30 bg-info-muted px-3 py-2 text-sm text-info">
      <Badge tone="running">running</Badge>
      A backup or restore job is in progress. Status refreshes automatically.
    </div>
  );
}
