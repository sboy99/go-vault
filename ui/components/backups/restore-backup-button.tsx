"use client";

import { useCallback, useState } from "react";
import { useRouter } from "next/navigation";
import { RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";
import { jobHref } from "@/lib/route-id";
import type { Job } from "@/lib/api/types";

type RestoreEnvelope = {
  success: boolean;
  data: { data: Job; demo: boolean } | null;
  error: string | null;
};

export function RestoreBackupButton({
  backupId,
  dbName,
}: {
  backupId: string;
  dbName: string;
}) {
  const router = useRouter();
  const { push } = useToast();
  const [open, setOpen] = useState(false);
  const [confirm, setConfirm] = useState("");
  const [pending, setPending] = useState(false);

  const onClose = useCallback(() => {
    if (pending) return;
    setOpen(false);
    setConfirm("");
  }, [pending]);

  return (
    <>
      <button
        type="button"
        aria-label="Restore"
        onClick={() => setOpen(true)}
        className="inline-flex size-8 items-center justify-center rounded-none border border-dashed border-brand bg-transparent text-brand transition-colors duration-200 ease-out hover:bg-brand/10"
      >
        <RotateCcw className="size-4" strokeWidth={1.75} />
      </button>
      <Dialog
        open={open}
        onClose={onClose}
        title="Restore backup"
        description={`Type the database name (${dbName || "configured name"}) to confirm. This overwrites the live database.`}
      >
        <label className="block text-[13px] text-foreground-muted">
          Database name
          <Input
            className="mt-1"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder={dbName || "database name"}
            autoComplete="off"
            disabled={pending}
          />
        </label>
        <div className="mt-4 flex items-center justify-end gap-3 border-t border-dashed border-edge pt-3">
          <Button variant="secondary" onClick={onClose} disabled={pending}>
            Cancel
          </Button>
          <Button
            variant="primary"
            className="rounded-none"
            disabled={pending || !confirm.trim()}
            onClick={async () => {
              setPending(true);
              try {
                const res = await fetch("/api/restores", {
                  method: "POST",
                  headers: { "Content-Type": "application/json" },
                  body: JSON.stringify({
                    backup_id: backupId,
                    confirm: confirm.trim(),
                  }),
                });
                const body = (await res.json()) as RestoreEnvelope;
                if (!res.ok || !body.success || !body.data?.data) {
                  push(body.error || "Restore failed", "error");
                  return;
                }
                push("Restore job started", "success");
                setOpen(false);
                setConfirm("");
                router.push(jobHref(body.data.data.id));
                router.refresh();
              } catch {
                push("Restore failed", "error");
              } finally {
                setPending(false);
              }
            }}
          >
            {pending ? "Starting…" : "Restore"}
          </Button>
        </div>
      </Dialog>
    </>
  );
}
