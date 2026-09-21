import { NextResponse } from "next/server";
import { ApiError, getJob } from "@/lib/api/client";

export async function GET(
  _req: Request,
  ctx: RouteContext<"/api/jobs/[id]">,
) {
  try {
    const { id } = await ctx.params;
    const result = await getJob(id);
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
