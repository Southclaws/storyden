import Image from "next/image";
import { useQueryState } from "nuqs";

import { Asset } from "@/api/openapi-schema";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { AssetLightbox } from "./AssetLightbox";

type Props = {
  asset: Asset;
  set?: Asset[];
  setIndex?: number;
  width?: number;
  height?: number;
};

const thumbnailStyles = css({
  borderRadius: "sm",
  height: "full",
  objectFit: "cover",
  overflowClipMargin: "unset",
});

export function AssetThumbnail({
  asset,
  set,
  setIndex,
  width = 80,
  height = 80,
}: Props) {
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
      cursor="pointer"
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

      <AssetLightbox
        asset={asset}
        set={set}
        setIndex={setIndex}
        present={lightbox}
        onClose={handleClose}
      />
    </Box>
  );
}
