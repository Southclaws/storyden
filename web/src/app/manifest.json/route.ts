import { Manifest } from "next/dist/lib/metadata/types/manifest-types";
import { NextResponse } from "next/server";

import { getInfo } from "src/utils/info";

export async function GET() {
  const info = await getInfo();

  const manifest: Manifest = {
    id: "/",
    name: info.title ?? "Storyden",
    short_name: info.title ?? "Storyden",
    description: info.description,
    display: "fullscreen",
    start_url: "/",
    // TODO: figure out a good choice for this.
    theme_color: "white",
    // TODO: make sure this makes sense.
    background_color: info.accent_colour,
    icons: [
      { src: "/api/v1/info/icon/512x512", sizes: "512x512", type: "image/png" },
      { src: "/api/v1/info/icon/32x32", sizes: "32x32", type: "image/png" },
      { src: "/api/v1/info/icon/180x180", sizes: "180x180", type: "image/png" },
      { src: "/api/v1/info/icon/120x120", sizes: "120x120", type: "image/png" },
      { src: "/api/v1/info/icon/167x167", sizes: "167x167", type: "image/png" },
      { src: "/api/v1/info/icon/152x152", sizes: "152x152", type: "image/png" },
    ],
  };

  return new NextResponse(JSON.stringify(manifest), {
    headers: {
      "Content-Type": "application/manifest+json",
    },
  });
}
