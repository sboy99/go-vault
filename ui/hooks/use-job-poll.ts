"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import type { Job } from "@/lib/api/types";

type JobsEnvelope = {
  success: boolean;
  data: { data: Job[]; demo: boolean } | null;
};

export function useJobPoll(enabled: boolean, intervalMs = 3000) {
  const router = useRouter();

  useEffect(() => {
    if (!enabled) return;

    let cancelled = false;
    let timer: number | undefined;

    const tick = async () => {
      try {
        const res = await fetch("/api/jobs?limit=10");
        const body = (await res.json()) as JobsEnvelope;
        const jobs = body.data?.data ?? [];
        const running = jobs.some(
          (j) => j.status === "running" || j.status === "queued",
        );
        if (!cancelled && running) {
          router.refresh();
        }
        if (!running) {
          return;
        }
      } catch {
        // ignore transient poll errors
      }
      if (!cancelled) {
        timer = window.setTimeout(tick, intervalMs);
      }
    };

    timer = window.setTimeout(tick, intervalMs);
    return () => {
      cancelled = true;
      if (timer) window.clearTimeout(timer);
    };
  }, [enabled, intervalMs, router]);
}
