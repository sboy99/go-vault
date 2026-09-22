import { ApiError, downloadBackupStream } from "@/lib/api/client";
import { fromPathId } from "@/lib/route-id";

export async function GET(
  _req: Request,
  ctx: { params: Promise<{ id: string }> },
) {
  try {
    const { id: pathId } = await ctx.params;
    return await downloadBackupStream(fromPathId(pathId));
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json(
        { success: false, data: null, error: err.message },
        { status: err.status },
      );
    }
    return Response.json(
      { success: false, data: null, error: "internal error" },
      { status: 500 },
    );
  }
}
