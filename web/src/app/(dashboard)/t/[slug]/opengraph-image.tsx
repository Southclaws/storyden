import { ImageResponse } from "next/og";

import { getInfo } from "src/utils/info";

import { threadGet } from "@/api/openapi-server/threads";
import { getAssetURL } from "@/utils/asset";

import { Props } from "./page";

export const size = {
  width: 1200,
  height: 630,
};

export const contentType = "image/png";

export default async function Image({ params }: Props) {
  const { slug } = await params;
  const { data } = await threadGet(slug);

  const { accent_colour } = await getInfo();

  const image = data.assets[0];

  return new ImageResponse(
    (
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "flex-start",
          justifyContent: "flex-start",
          backgroundColor: accent_colour,
          fontSize: 16,
        }}
      >
        {image ? (
          // eslint-disable-next-line @next/next/no-img-element, jsx-a11y/alt-text
          (<img
            src={getAssetURL(image.path)}
            width="100%"
            height="100%"
            style={{
              objectPosition: "center",
              objectFit: "cover",
            }}
          />)
        ) : (
          <div></div>
        )}

        <div
          style={{
            position: "absolute",
            bottom: 0,
            display: "flex",
            flexDirection: "column",
            padding: "2rem",
            backgroundColor: "hsla(180deg, 10%, 10%, 0.58)",
            color: "white",
          }}
        >
          <div
            style={{
              fontSize: "4rem",
              fontWeight: 600,
            }}
          >
            {data.title}
          </div>

          <div
            style={{
              fontSize: "2rem",
              fontWeight: 300,
            }}
          >
            {data.description}
          </div>
        </div>
      </div>
    ),
    {
      ...size,
      emoji: "fluent",
    },
  );
}
