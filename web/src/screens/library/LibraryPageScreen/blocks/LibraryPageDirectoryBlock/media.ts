import type { Asset } from "@/api/openapi-schema";
import { getAssetImageProps } from "@/utils/asset";

export function getDirectoryImageProps(
  image: Pick<Asset, "path" | "width" | "height"> | undefined,
  layout: "grid" | "gallery",
) {
  const { src, srcSet } = getAssetImageProps(image);
  const objectFit: "contain" | "cover" =
    layout === "gallery" &&
    image &&
    image.width > 0 &&
    image.height > image.width * 3
      ? "contain"
      : "cover";

  return {
    src,
    srcSet,
    // A second 200px column fits once the directory reaches 408px.
    sizes: "408px",
    htmlWidth: image && image.width > 0 ? image.width : undefined,
    htmlHeight: image && image.height > 0 ? image.height : undefined,
    objectFit,
  };
}
