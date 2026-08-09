import Image from "next/image";

import { BannerEditor } from "@/components/admin/BrandSettings/BannerEditor/BannerEditor";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getBannerURL } from "@/utils/icon";

import { useFeedBlockEditor } from "../Context";

export function FeedCoverBlock() {
  const { isEditing } = useFeedBlockEditor();

  if (isEditing) {
    return <BannerEditor />;
  }

  return (
    <Box width="full" height="64" position="relative">
      <Image
        className={css({
          width: "full",
          height: "full",
          borderRadius: "md",
          objectFit: "cover",
        })}
        src={getBannerURL()}
        alt=""
        fill
        sizes="(max-width: 768px) 100vw, 960px"
      />
    </Box>
  );
}
