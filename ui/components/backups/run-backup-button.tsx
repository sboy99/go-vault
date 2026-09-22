"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/toast";
import { jobHref } from "@/lib/route-id";
import type { Job } from "@/lib/api/types";

type CreateEnvelope = {
  success: boolean;
  data: { data: Job; demo: boolean } | null;
  error: string | null;
};

export function RunBackupButton() {
  const router = useRouter();
  const { push } = useToast();
  const [pending, setPending] = useState(false);

  return (
    <Button
      variant="secondary"
      disabled={pending}
      onClick={async () => {
        setPending(true);
        try {
          const res = await fetch("/api/backups", { method: "POST" });
          const body = (await res.json()) as CreateEnvelope;
          if (!res.ok || !body.success || !body.data?.data) {
            push(body.error || "Failed to start backup", "error");
            return;
          }
          push("Backup job started", "success");
          router.push(jobHref(body.data.data.id));
          router.refresh();
        } catch {
          push("Failed to start backup", "error");
        } finally {
          setPending(false);
        }
      }}
    >
      {pending ? "Starting…" : "Run backup"}
    </Button>
  );
}
