"use client";

import { useJobPoll } from "@/hooks/use-job-poll";

export function JobsPoller({ enabled }: { enabled: boolean }) {
  useJobPoll(enabled);
  return null;
}
