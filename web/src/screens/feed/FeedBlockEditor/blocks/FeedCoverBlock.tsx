import Image from "next/image";

import { FeedBlock } from "@/lib/settings/feed";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

export function FeedCoverBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "cover" }>;
}) {
  const src = getAssetURL(block.asset?.path, "l");
  if (!src) {
    return null;
  }

  return (
    <Box
      className="feed-cover"
      width="full"
      aspectRatio={4}
      position="relative"
    >
      <Image
        className={
          "feed-cover__image " +
          css({
            width: "full",
            height: "full",
            borderRadius: "md",
            objectFit: "cover",
          })
        }
        src={src}
        alt=""
        loading="eager"
        fill
        sizes="(max-width: 768px) 100vw, 960px"
      />
    </Box>
  );
}
