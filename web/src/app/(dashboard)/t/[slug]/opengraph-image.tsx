import { ImageResponse } from "next/server";

import { server } from "src/api/client";
import { ThreadGetResponse } from "src/api/openapi/schemas";
import { getInfo } from "src/utils/info";

import { Props } from "./page";

export const size = {
  width: 1200,
  height: 630,
};

export const contentType = "image/png";

export default async function Image({ params: { slug } }: Props) {
  const thread = await server<ThreadGetResponse>({
    url: `/v1/threads/${slug}`,
  });

  const { accent_colour } = await getInfo();

  const image = thread.assets[0];

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
          <img
            src={image.url}
            width="100%"
            height="100%"
            style={{
              objectPosition: "center",
              objectFit: "cover",
            }}
          />
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
            {thread.title}
          </div>

          <div
            style={{
              fontSize: "2rem",
              fontWeight: 300,
            }}
          >
            {thread.short}
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
