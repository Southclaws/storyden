import { ImageResponse } from "next/og";

import { API_ADDRESS } from "src/config";

export const dynamic = "force-dynamic";

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
        src={`${API_ADDRESS}/api/v1/info/icon/512x512`}
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
