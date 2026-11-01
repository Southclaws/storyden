import type { Asset } from "@/api/openapi-schema";
import type { CardItem } from "@/components/ui/surface";
import { getAssetImageProps, getAssetURL } from "@/utils/asset";

import {
  DirectoryCardLayout,
  RichCardShape,
  getRichCardImageAspectRatio,
} from "./directoryCard";

export function getDirectoryImageProps(
  image: Pick<Asset, "path" | "width" | "height"> | undefined,
  layout: "grid" | "gallery",
) {
  const { src, srcSet } = getAssetImageProps(image);
  let objectFit: "contain" | "cover" = "cover";
  let htmlWidth: number | undefined;
  let htmlHeight: number | undefined;

  if (image) {
    if (image.width > 0) {
      htmlWidth = image.width;
    }

    if (image.height > 0) {
      htmlHeight = image.height;
    }

    if (layout === "gallery" && image.width > 0) {
      if (image.height > image.width * 3) {
        objectFit = "contain";
      }
    }
  }

  return {
    src,
    srcSet,
    // A second 200px column fits once the directory reaches 408px.
    sizes: "408px",
    htmlWidth,
    htmlHeight,
    objectFit,
  };
}

export function getDirectoryCardImage(
  image: Asset | undefined,
  layout: DirectoryCardLayout,
  shape: RichCardShape,
): Pick<CardItem, "image" | "imageAspectRatio" | "imageFit" | "imageProps"> {
  const imageAspectRatio = getRichCardImageAspectRatio(layout, shape, image);

  if (shape !== "minimal") {
    return { image: getAssetURL(image?.path, "m"), imageAspectRatio };
  }

  let imageLayout: "grid" | "gallery" = "grid";

  if (layout === "gallery") {
    imageLayout = "gallery";
  }

  const media = getDirectoryImageProps(image, imageLayout);

  return {
    image: media.src,
    imageAspectRatio,
    imageFit: media.objectFit,
    imageProps: {
      srcSet: media.srcSet,
      sizes: media.sizes,
      width: media.htmlWidth,
      height: media.htmlHeight,
      loading: "lazy",
      decoding: "async",
    },
  };
}
