import "server-only";

import {
  mockBackupById,
  mockJobById,
  mocks,
} from "@/lib/mock/fixtures";
import type {
  Backup,
  DataResult,
  Envelope,
  HealthStatus,
  Job,
  RedactedConfig,
  RetentionPreview,
  StatsResponse,
} from "@/lib/api/types";

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

function baseUrl(): string {
  return (process.env.GOVAULT_API_URL || "http://127.0.0.1:8080").replace(/\/$/, "");
}

function forceMocks(): boolean {
  return process.env.GOVAULT_USE_MOCKS === "1";
}

async function apiFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const headers = new Headers(init?.headers);
  headers.set("Accept", "application/json");
  if (init?.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  const token = process.env.GOVAULT_API_TOKEN;
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${baseUrl()}${path}`, {
    ...init,
    headers,
    cache: "no-store",
  });

  const contentType = res.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    if (!res.ok) {
      throw new ApiError(res.statusText || "request failed", res.status);
    }
    return undefined as T;
  }

  const envelope = (await res.json()) as Envelope<T>;
  if (!res.ok || !envelope.success) {
    throw new ApiError(envelope.error || res.statusText || "request failed", res.status);
  }
  return envelope.data as T;
}

function isConnectionError(err: unknown): boolean {
  if (err instanceof TypeError) return true;
  if (err instanceof Error) {
    const msg = err.message.toLowerCase();
    return (
      msg.includes("fetch failed") ||
      msg.includes("econnrefused") ||
      msg.includes("enotfound") ||
      msg.includes("network")
    );
  }
  return false;
}

function shouldUseFallback(err: unknown): boolean {
  if (isConnectionError(err)) return true;
  // Old servers without /v1/stats|/v1/config|/v1/retention return 404 HTML.
  if (err instanceof ApiError && (err.status === 404 || err.status === 502 || err.status === 503)) {
    return true;
  }
  return false;
}

async function withFallback<T>(
  live: () => Promise<T>,
  fallback: T,
): Promise<DataResult<T>> {
  if (forceMocks()) {
    return { data: fallback, demo: true };
  }
  try {
    const data = await live();
    return { data, demo: false };
  } catch (err) {
    if (shouldUseFallback(err)) {
      return { data: fallback, demo: true };
    }
    throw err;
  }
}

export async function listBackups(opts?: {
  limit?: number;
  offset?: number;
  status?: string;
}): Promise<DataResult<Backup[]>> {
  const params = new URLSearchParams();
  if (opts?.limit) params.set("limit", String(opts.limit));
  if (opts?.offset) params.set("offset", String(opts.offset));
  if (opts?.status) params.set("status", opts.status);
  const qs = params.toString();
  return withFallback(
    () => apiFetch<Backup[]>(`/v1/backups${qs ? `?${qs}` : ""}`),
    mocks.backups.filter((b) => !opts?.status || b.status === opts.status),
  );
}

export async function getBackup(id: string): Promise<DataResult<Backup>> {
  return withFallback(async () => {
    const data = await apiFetch<Backup>(`/v1/backups/${encodeURIComponent(id)}`);
    return data;
  }, mockBackupById(id) ?? mocks.backups[0]);
}

export async function createBackup(): Promise<DataResult<Job>> {
  if (forceMocks()) {
    return {
      data: {
        id: `job-mock-${Date.now()}`,
        type: "backup",
        status: "running",
        created_at: new Date().toISOString(),
        started_at: new Date().toISOString(),
      },
      demo: true,
    };
  }
  try {
    const data = await apiFetch<Job>("/v1/backups", { method: "POST" });
    return { data, demo: false };
  } catch (err) {
    if (shouldUseFallback(err)) {
      return {
        data: {
          id: `job-mock-${Date.now()}`,
          type: "backup",
          status: "running",
          created_at: new Date().toISOString(),
          started_at: new Date().toISOString(),
        },
        demo: true,
      };
    }
    throw err;
  }
}

export async function restoreBackup(
  backupId: string,
  confirm: string,
): Promise<DataResult<Job>> {
  if (forceMocks()) {
    return {
      data: {
        id: `job-mock-restore-${Date.now()}`,
        type: "restore",
        status: "running",
        backup_id: backupId,
        created_at: new Date().toISOString(),
        started_at: new Date().toISOString(),
      },
      demo: true,
    };
  }
  try {
    const data = await apiFetch<Job>("/v1/restores", {
      method: "POST",
      body: JSON.stringify({ backup_id: backupId, confirm }),
    });
    return { data, demo: false };
  } catch (err) {
    if (shouldUseFallback(err)) {
      return {
        data: {
          id: `job-mock-restore-${Date.now()}`,
          type: "restore",
          status: "running",
          backup_id: backupId,
          created_at: new Date().toISOString(),
          started_at: new Date().toISOString(),
        },
        demo: true,
      };
    }
    throw err;
  }
}

export async function listJobs(opts?: {
  limit?: number;
  offset?: number;
}): Promise<DataResult<Job[]>> {
  const params = new URLSearchParams();
  if (opts?.limit) params.set("limit", String(opts.limit));
  if (opts?.offset) params.set("offset", String(opts.offset));
  const qs = params.toString();
  return withFallback(
    () => apiFetch<Job[]>(`/v1/jobs${qs ? `?${qs}` : ""}`),
    mocks.jobs,
  );
}

export async function getJob(id: string): Promise<DataResult<Job>> {
  return withFallback(
    () => apiFetch<Job>(`/v1/jobs/${encodeURIComponent(id)}`),
    mockJobById(id) ?? mocks.jobs[0],
  );
}

export async function getStats(): Promise<DataResult<StatsResponse>> {
  return withFallback(() => apiFetch<StatsResponse>("/v1/stats"), mocks.stats);
}

export async function getConfig(): Promise<DataResult<RedactedConfig>> {
  return withFallback(() => apiFetch<RedactedConfig>("/v1/config"), mocks.config);
}

export async function getRetention(): Promise<DataResult<RetentionPreview>> {
  return withFallback(
    () => apiFetch<RetentionPreview>("/v1/retention"),
    mocks.retention,
  );
}

export async function getHealth(): Promise<DataResult<HealthStatus>> {
  if (forceMocks()) {
    return { data: mocks.health, demo: true };
  }
  try {
    const health = await probe("/healthz");
    const ready = await probe("/readyz");
    return {
      data: {
        healthz: health.ok,
        readyz: ready.ok,
        health_error: health.error,
        ready_error: ready.error,
      },
      demo: false,
    };
  } catch {
    return { data: mocks.health, demo: true };
  }
}

async function probe(path: string): Promise<{ ok: boolean; error?: string }> {
  const headers = new Headers({ Accept: "application/json" });
  const token = process.env.GOVAULT_API_TOKEN;
  if (token) headers.set("Authorization", `Bearer ${token}`);
  const res = await fetch(`${baseUrl()}${path}`, { headers, cache: "no-store" });
  if (!res.ok) {
    try {
      const body = (await res.json()) as Envelope<unknown>;
      return { ok: false, error: body.error || res.statusText };
    } catch {
      return { ok: false, error: res.statusText };
    }
  }
  return { ok: true };
}

export async function downloadBackupStream(id: string): Promise<Response> {
  if (forceMocks()) {
    return new Response("demo backup payload", {
      headers: {
        "Content-Type": "application/octet-stream",
        "Content-Disposition": `attachment; filename="demo-backup.dump"`,
      },
    });
  }

  const headers = new Headers();
  const token = process.env.GOVAULT_API_TOKEN;
  if (token) headers.set("Authorization", `Bearer ${token}`);

  try {
    const res = await fetch(
      `${baseUrl()}/v1/backups/${encodeURIComponent(id)}/download`,
      { headers, cache: "no-store" },
    );
    if (!res.ok) {
      throw new ApiError("download failed", res.status);
    }
    return new Response(res.body, {
      status: res.status,
      headers: {
        "Content-Type":
          res.headers.get("Content-Type") || "application/octet-stream",
        "Content-Disposition":
          res.headers.get("Content-Disposition") ||
          `attachment; filename="backup.dump"`,
      },
    });
  } catch (err) {
    if (shouldUseFallback(err)) {
      return new Response("demo backup payload", {
        headers: {
          "Content-Type": "application/octet-stream",
          "Content-Disposition": `attachment; filename="demo-backup.dump"`,
        },
      });
    }
    throw err;
  }
}
