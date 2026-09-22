import { NextResponse } from "next/server";
import { ApiError, getJob } from "@/lib/api/client";
import { fromPathId } from "@/lib/route-id";

export async function GET(
  _req: Request,
  ctx: { params: Promise<{ id: string }> },
) {
  try {
    const { id: pathId } = await ctx.params;
    const result = await getJob(fromPathId(pathId));
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
