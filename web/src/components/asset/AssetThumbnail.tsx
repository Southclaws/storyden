import Image from "next/image";
import { useQueryState } from "nuqs";

import { Asset } from "@/api/openapi-schema";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { AssetLightbox } from "./AssetLightbox";

type Props = {
  asset: Asset;
  width?: number;
  height?: number;
};

const thumbnailStyles = css({
  borderRadius: "sm",
});

export function AssetThumbnail({ asset, width = 64, height = 64 }: Props) {
  const [view, setView] = useQueryState<string | null>("view", {
    defaultValue: null,
    clearOnDefault: true,
    parse: (value) => (value === "" ? null : value),
  });
  const url = getAssetURL(asset.path)!;

  const lightbox = view === asset.id;

  function handleOpen() {
    setView(asset.id);
  }

  function handleClose() {
    setView(null);
  }

  return (
    <Box
      w="20"
      h="20"
      borderRadius="md"
      overflow="hidden"
      grayscale="0.8"
      filter={lightbox ? "auto" : undefined}
    >
      <Image
        className={thumbnailStyles}
        src={url}
        alt={asset.filename}
        width={width}
        height={height}
        onClick={handleOpen}
      />

      <AssetLightbox asset={asset} present={lightbox} onClose={handleClose} />
    </Box>
  );
}
