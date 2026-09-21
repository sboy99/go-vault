import { NextResponse } from "next/server";
import { ApiError, listJobs } from "@/lib/api/client";

export async function GET(req: Request) {
  try {
    const url = new URL(req.url);
    const result = await listJobs({
      limit: Number(url.searchParams.get("limit") || 50),
      offset: Number(url.searchParams.get("offset") || 0),
    });
    return NextResponse.json({ success: true, data: result, error: null });
  } catch (err) {
    if (err instanceof ApiError) {
      return NextResponse.json(
        { success: false, data: null, error: err.message },
        { status: err.status },
      );
    }
    return NextResponse.json(
      { success: false, data: null, error: "internal error" },
      { status: 500 },
    );
  }
}
