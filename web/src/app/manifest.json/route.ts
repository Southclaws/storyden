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
      {
        src: "/icon-192x192.png",
        sizes: "192x192",
        type: "image/png",
      },
      {
        src: "/icon-256x256.png",
        sizes: "256x256",
        type: "image/png",
      },
      {
        src: "/icon-384x384.png",
        sizes: "384x384",
        type: "image/png",
      },
      {
        src: "/icon-512x512.png",
        sizes: "512x512",
        type: "image/png",
      },
    ],
  };

  return new NextResponse(JSON.stringify(manifest), {
    headers: {
      "Content-Type": "application/manifest+json",
    },
  });
}
