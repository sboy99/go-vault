import { NextResponse } from "next/server";
import { ApiError, createBackup } from "@/lib/api/client";

export async function POST() {
  try {
    const result = await createBackup();
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
