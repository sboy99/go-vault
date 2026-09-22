import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { toPathId } from "@/lib/route-id";

/**
 * Soft gate for BFF mutations: reject cross-site POSTs.
 * The go-vault API bearer token stays server-side only.
 * Bind this UI to localhost / a private network — it proxies privileged ops.
 *
 * Also redirects /jobs and /backups paths that still contain `:` to the
 * path-safe `~` form so old bookmarks cannot start an RSC prefetch livelock.
 */
export function proxy(req: NextRequest) {
  const { pathname } = req.nextUrl;
  const colonRedirect = redirectColonPath(pathname);
  if (colonRedirect) {
    const url = req.nextUrl.clone();
    url.pathname = colonRedirect;
    return NextResponse.redirect(url);
  }

  if (req.method === "GET" || req.method === "HEAD" || req.method === "OPTIONS") {
    return NextResponse.next();
  }

  if (!pathname.startsWith("/api/")) {
    return NextResponse.next();
  }

  const site = req.headers.get("sec-fetch-site");
  if (site && site !== "same-origin" && site !== "none") {
    return NextResponse.json(
      { success: false, data: null, error: "cross-site request blocked" },
      { status: 403 },
    );
  }

  return NextResponse.next();
}

function redirectColonPath(pathname: string): string | null {
  const match = pathname.match(/^\/(jobs|backups)\/(.+)$/);
  if (!match) return null;
  const [, kind, rawId] = match;
  if (!rawId.includes(":")) return null;
  // Only rewrite the id segment; leave trailing subpaths (e.g. /download) alone.
  const slash = rawId.indexOf("/");
  const idPart = slash === -1 ? rawId : rawId.slice(0, slash);
  const rest = slash === -1 ? "" : rawId.slice(slash);
  if (!idPart.includes(":")) return null;
  return `/${kind}/${toPathId(idPart)}${rest}`;
}

export const config = {
  matcher: ["/api/:path*", "/jobs/:path*", "/backups/:path*"],
};
