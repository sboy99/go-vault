import { NextResponse } from "next/server";
import { ApiError, restoreBackup } from "@/lib/api/client";

export async function POST(req: Request) {
  try {
    const body = (await req.json()) as { backup_id?: string; confirm?: string };
    if (!body.backup_id || !body.confirm) {
      return NextResponse.json(
        { success: false, data: null, error: "backup_id and confirm are required" },
        { status: 400 },
      );
    }
    const result = await restoreBackup(body.backup_id, body.confirm);
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
