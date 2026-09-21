import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

/**
 * Soft gate for BFF mutations: reject cross-site POSTs.
 * The go-vault API bearer token stays server-side only.
 * Bind this UI to localhost / a private network — it proxies privileged ops.
 */
export function proxy(req: NextRequest) {
  if (req.method === "GET" || req.method === "HEAD" || req.method === "OPTIONS") {
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

export const config = {
  matcher: ["/api/:path*"],
};
