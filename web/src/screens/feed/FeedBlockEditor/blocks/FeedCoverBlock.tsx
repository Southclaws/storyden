import Image from "next/image";

import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getBannerURL } from "@/utils/icon";

export function FeedCoverBlock() {
  return (
    <Box className="feed-cover" width="full" height="64" position="relative">
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
        src={getBannerURL()}
        alt=""
        fill
        sizes="(max-width: 768px) 100vw, 960px"
      />
    </Box>
  );
}
