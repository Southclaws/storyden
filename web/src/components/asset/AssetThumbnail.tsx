import { useQueryState } from "nuqs";

import { Asset } from "@/api/openapi-schema";
import { css } from "@/styled-system/css";
import { Box, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { IconButton } from "../ui/icon-button";
import { DeleteIcon } from "../ui/icons/Delete";

import { AssetLightbox } from "./AssetLightbox";

type Props = {
  className?: string;
  asset: Asset;
  set?: Asset[];
  setIndex?: number;
  showDeleteButton?: boolean;
  handleDelete?: () => Promise<void>;
};

const thumbnailStyles = css({
  borderRadius: "sm",
  width: "full",
  height: "full",
  objectFit: "cover",
  overflowClipMargin: "unset",
});

export function AssetThumbnail({
  className,
  asset,
  set,
  setIndex,
  showDeleteButton = false,
  handleDelete = undefined,
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
      w="full"
      h="full"
      borderRadius="md"
      overflow="hidden"
      grayscale="0.8"
      cursor="pointer"
      filter={lightbox ? "auto" : undefined}
    >
      <styled.img
        className={className ?? thumbnailStyles}
        src={url}
        alt={asset.filename}
        onClick={handleOpen}
        aspectRatio={`1`}
      />

      <AssetLightbox
        asset={asset}
        set={set}
        setIndex={setIndex}
        present={lightbox}
        onClose={handleClose}
      />

      {showDeleteButton && (
        <IconButton
          type="button"
          position="absolute"
          top="1"
          right="1"
          colorPalette="tomato"
          variant="subtle"
          size="xs"
          w="5"
          h="5"
          minW="5"
          title="Remove media from page"
          onClick={handleDelete}
        >
          <DeleteIcon />
        </IconButton>
      )}
    </Box>
  );
}
