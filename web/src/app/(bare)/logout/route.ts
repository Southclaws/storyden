import { NextResponse } from "next/server";

import { WEB_ADDRESS } from "src/config";

export async function GET() {
  return NextResponse.redirect(WEB_ADDRESS, {
    headers: {
      "Clear-Site-Data": `"*"`,
      "Set-Cookie": "storyden-session=",
    },
  });
}
