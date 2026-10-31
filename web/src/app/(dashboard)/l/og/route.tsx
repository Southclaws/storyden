import { ImageResponse } from "next/og";
import { NextRequest, NextResponse } from "next/server";

import { nodeGet } from "@/api/openapi-server/nodes";
import { loadMetadataImage } from "@/lib/metadata/image-response";
import { getSettings } from "@/lib/settings/settings-server";

const size = {
  width: 1200,
  height: 630,
};

export async function GET(req: NextRequest) {
  const slug = req.nextUrl.searchParams.get("slug");
  if (!slug) {
    return new NextResponse("No slug provided", { status: 400 });
  }

  const { data } = await nodeGet(slug, undefined, {
    anonymous: true,
    cache: "no-store",
    signal: req.signal,
  });
  if (data.visibility !== "published") {
    return new NextResponse("Not Found", {
      status: 404,
      headers: { "Cache-Control": "no-store" },
    });
  }

  const { accent_colour } = await getSettings();

  const image =
    data.primary_image?.parent ??
    data.link?.primary_image ??
    data.primary_image;

  const imageData = image
    ? await loadMetadataImage(
        `/assets/${encodeURIComponent(image.filename)}?variant=l`,
      )
    : undefined;
  if (imageData instanceof Response) return imageData;

  return new ImageResponse(
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
      {imageData ? (
        // eslint-disable-next-line @next/next/no-img-element, jsx-a11y/alt-text
        <img
          // Satori accepts binary image sources; React's DOM types only allow URLs.
          // @ts-expect-error ArrayBuffer is supported by next/og.
          src={imageData}
          width="100%"
          height="100%"
          style={{
            objectPosition: "center",
            objectFit: "cover",
          }}
        />
      ) : (
        <div
          style={{
            position: "absolute",
            bottom: 0,
            display: "flex",
            flexDirection: "column",
            padding: "2rem",
            width: "100%",
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
            {data.name}
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
      )}
    </div>,
    {
      ...size,
      emoji: "fluent",
    },
  );
}
