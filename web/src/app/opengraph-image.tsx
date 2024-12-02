/* eslint-disable jsx-a11y/alt-text, @next/next/no-img-element */
import { ImageResponse } from "next/og";

import { getSettings } from "@/lib/settings/settings-server";
import { getBannerURL, getIconURL } from "@/utils/icon";

export const size = {
  width: 1200,
  height: 630,
};

export const runtime = "edge";
export const contentType = "image/png";

export default async function Image() {
  const settings = await getSettings();
  const iconURL = getIconURL("512x512");
  const backgroundImageURL = getBannerURL();

  const { title, accent_colour } = settings;

  return new ImageResponse(
    (
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#fff",
        }}
      >
        <img
          src={backgroundImageURL}
          width="100%"
          height="100%"
          style={{
            objectPosition: "center",
            objectFit: "cover",
          }}
        />
      </div>
    ),
    {
      ...size,
      emoji: "fluent",
    },
  );
}
