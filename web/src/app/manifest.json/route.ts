import { Manifest } from "next/dist/lib/metadata/types/manifest-types";
import { NextResponse } from "next/server";

import { GetInfoOKResponse } from "src/api/openapi/schemas";
import { API_ADDRESS } from "src/config";

export async function GET() {
  const res = await fetch(`${API_ADDRESS}/api/v1/info`);
  if (!res.ok) {
    throw new Error(
      `failed to fetch API info endpoint: ${res.status} ${res.statusText}`
    );
  }

  const info = (await res.json()) as GetInfoOKResponse;

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
