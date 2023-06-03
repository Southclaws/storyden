import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const url = request.nextUrl.clone();
  url.pathname = "/";

  return NextResponse.redirect(url, {
    headers: {
      "Clear-Site-Data": `"*"`,
    },
  });
}
