import { NextResponse } from "next/server";
import { getHealth } from "@/lib/api/client";

export async function GET() {
  const result = await getHealth();
  return NextResponse.json({ success: true, data: result, error: null });
}
