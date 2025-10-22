import { ImageResponse } from "next/og";

import { getIconURL } from "@/utils/icon";

export const size = {
  width: 512,
  height: 512,
};
export const contentType = "image/png";

export default function Icon() {
  return new ImageResponse(
    (
      // eslint-disable-next-line @next/next/no-img-element, jsx-a11y/alt-text
      <img
        src={getIconURL("512x512")}
        width={512}
        height={512}
        sizes="512x512"
      />
    ),
    {
      ...size,
    },
  );
}
