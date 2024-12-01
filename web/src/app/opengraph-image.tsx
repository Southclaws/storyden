/* eslint-disable jsx-a11y/alt-text */

/* eslint-disable @next/next/no-img-element */
import { ImageResponse } from "next/og";

import { interBold, interRegular } from "@/app/fonts/og";
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
          fontFamily: 'Inter, "Material Icons"',
          fontSize: 32,
          fontWeight: 600,
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

        <div
          style={{
            position: "absolute",
            top: 0,
            left: 0,
            display: "flex",
            flexDirection: "column",
            width: "100%",
            padding: "2rem",
            color: "white",
          }}
        >
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              gap: "2rem",

              background: "hsla(180deg, 10%, 10%, 58%)",
              borderRadius: "12px",
              padding: "1rem",
            }}
          >
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                fontSize: "4rem",
                fontWeight: 600,
              }}
            >
              <span>{title}</span>
            </div>

            <img
              src={iconURL}
              width="100"
              height="100"
              style={{
                objectPosition: "center",
                objectFit: "cover",
                borderRadius: "12px",
              }}
            />
          </div>
        </div>
      </div>
    ),
    {
      ...size,
      fonts: [
        {
          name: "Inter",
          data: await interRegular(),
          style: "normal",
          weight: 400,
        },
        {
          name: "Inter",
          data: await interBold(),
          style: "normal",
          weight: 800,
        },
      ],
      emoji: "fluent",
    },
  );
}
