import { revalidatePath, revalidateTag } from "next/cache";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

import { WEB_ADDRESS } from "src/config";

const cookieName = "storyden-session";

export async function POST() {
  (await cookies()).delete(cookieName);

  revalidateTag("accounts", { expire: 0 });
  revalidatePath("/", "layout");

  return NextResponse.redirect(WEB_ADDRESS, {
    headers: {
      "Clear-Site-Data": `"*"`,
      "Cache-Control": "no-cache, no-store, must-revalidate",
    },
  });
}
