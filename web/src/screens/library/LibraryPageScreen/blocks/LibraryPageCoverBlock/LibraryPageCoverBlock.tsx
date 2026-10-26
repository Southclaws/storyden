import Image from "next/image";
import { ReactNode, Suspense, lazy } from "react";

import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetImageProps, getAssetURL } from "@/utils/asset";

import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";
import { useBlock } from "../useBlock";

import { CROP_STENCIL_HEIGHT, CROP_STENCIL_WIDTH } from "./constants";
import { useLibraryPageCoverBlock } from "./useLibraryPageCoverBlock";

const LibraryPageCoverBlockEditing = lazy(() =>
  import("./LibraryPageCoverBlockEditing").then((module) => ({
    default: module.LibraryPageCoverBlockEditing,
  })),
);

const frameStyles = {
  banner: css({ width: "full", aspectRatio: "4/1" }),
  image: css({ width: "full" }),
  artwork: css({ width: "[min(100%,20rem)]", aspectRatio: "1" }),
};
const imageBase = css.raw({
  display: "block",
  mx: "auto",
  maxWidth: "full",
  objectPosition: "center",
});
const imageStyles = {
  banner: css(imageBase, { width: "full", height: "full", objectFit: "cover" }),
  image: css(imageBase, {
    width: "auto",
    height: "auto",
    maxHeight: "[min(70svh,40rem)]",
    objectFit: "contain",
  }),
  artwork: css(imageBase, {
    width: "full",
    height: "full",
    objectFit: "cover",
  }),
};

export function LibraryPageCoverBlock() {
  const { isDirectEditing } = useEditState();
  const primary_image = useWatch((s) => s.draft.primary_image);
  const block = useBlock("cover");
  const layout = block?.config?.layout ?? "banner";
  const imageProps = getAssetImageProps(primary_image, "l");

  // Resize variants round very thin dimensions to one pixel, changing their
  // proportions. These tiny originals are cheaper and preserve the true ratio.
  const useOriginal =
    primary_image?.width &&
    primary_image?.height &&
    Math.min(primary_image.width, primary_image.height) < 8;
  const src = useOriginal ? getAssetURL(primary_image.path) : imageProps.src;

  const cover = src ? (
    <Box
      data-main-image-layout={layout}
      className={frameStyles[layout]}
      mx="auto"
      overflow="hidden"
      borderRadius="md"
    >
      <picture
        className={css({
          display: "block",
          width: "full",
          height: "full",
        })}
      >
        <source
          srcSet={useOriginal ? undefined : imageProps.srcSet}
          sizes="(max-width: 960px) 100vw, 1200px"
        />
        <Image
          className={imageStyles[layout]}
          src={src}
          sizes="(max-width: 960px) 100vw, 1200px"
          loading="eager"
          alt={block?.config?.alt ?? ""}
          width={primary_image?.width || CROP_STENCIL_WIDTH}
          height={primary_image?.height || CROP_STENCIL_HEIGHT}
        />
      </picture>
    </Box>
  ) : null;

  return isDirectEditing ? <CoverEditingHost fallback={cover} /> : cover;
}

function CoverEditingHost({ fallback }: { fallback: ReactNode }) {
  // Subscribe before the cropper loads so link-block cover updates aren't lost.
  const editor = useLibraryPageCoverBlock();

  return (
    <Suspense fallback={fallback}>
      <LibraryPageCoverBlockEditing {...editor} />
    </Suspense>
  );
}
