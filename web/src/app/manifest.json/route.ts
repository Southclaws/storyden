import { Manifest } from "next/dist/lib/metadata/types/manifest-types";
import { NextResponse } from "next/server";

import { getColourAsHex } from "src/utils/colour";

import { API_ADDRESS } from "@/config";
import { getSettings } from "@/lib/settings/settings-server";

export async function GET() {
  const settings = await getSettings();

  const backgroundColour = getColourAsHex(settings.accent_colour);

  const manifest: Manifest = {
    id: "/",
    name: settings.title ?? "Storyden",
    short_name: settings.title ?? "Storyden",
    description: settings.description,
    display: "fullscreen",
    start_url: "/",
    // TODO: figure out a good choice for this.
    theme_color: "white",
    // TODO: make sure this makes sense.
    background_color: backgroundColour,

    // prettier-ignore
    icons: [
      { src: `${API_ADDRESS}/api/info/icon/512x512`, sizes: "512x512", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/32x32`, sizes: "32x32", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/180x180`, sizes: "180x180", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/120x120`, sizes: "120x120", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/167x167`, sizes: "167x167", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/152x152`, sizes: "152x152", type: "image/png" },
      { src: `${API_ADDRESS}/api/info/icon/512x512`, sizes: "512x512", type: "image/png", purpose: "maskable" },
      { src: `${API_ADDRESS}/api/info/icon/32x32`, sizes: "32x32", type: "image/png", purpose: "maskable" },
      { src: `${API_ADDRESS}/api/info/icon/180x180`, sizes: "180x180", type: "image/png", purpose: "maskable" },
      { src: `${API_ADDRESS}/api/info/icon/120x120`, sizes: "120x120", type: "image/png", purpose: "maskable" },
      { src: `${API_ADDRESS}/api/info/icon/167x167`, sizes: "167x167", type: "image/png", purpose: "maskable" },
      { src: `${API_ADDRESS}/api/info/icon/152x152`, sizes: "152x152", type: "image/png", purpose: "maskable" },
    ],
  };

  return new NextResponse(JSON.stringify(manifest), {
    headers: {
      "Content-Type": "application/manifest+json",
    },
  });
}
