/**
 * Job and backup ids are ISO timestamps (`2026-09-22T03:46:57…Z_hex`).
 * A raw `:` in a Next.js App Router path segment mismatches the RSC
 * cache key and can livelock prefetch. Swap `:` ↔ `~` for path use;
 * `~` is left alone by the URL parser and encodeURIComponent.
 */
export function toPathId(id: string): string {
  return id.replaceAll(":", "~");
}

export function fromPathId(pathId: string): string {
  return pathId.replaceAll("~", ":");
}

export function jobHref(id: string): string {
  return `/jobs/${toPathId(id)}`;
}

export function backupHref(id: string): string {
  return `/backups/${toPathId(id)}`;
}
