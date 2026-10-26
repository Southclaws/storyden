import { preload } from "react-dom";

import type { Asset } from "@/api/openapi-schema";
import { HTMLStyledProps, styled } from "@/styled-system/jsx";
import { getAssetImageProps } from "@/utils/asset";

type Props = Omit<
  HTMLStyledProps<"img">,
  "src" | "srcSet" | "sizes" | "htmlWidth" | "htmlHeight"
> & {
  asset: Pick<Asset, "path" | "width" | "height">;
  variant?: "s" | "m" | "l";
  sizes: string;
  priority?: boolean;
};

export function AssetImage({
  asset,
  variant = "m",
  sizes,
  priority = false,
  alt = "",
  ...rest
}: Props) {
  const { src, srcSet, width, height } = getAssetImageProps(asset, variant);
  const responsiveSizes = srcSet ? sizes : undefined;

  if (priority && src) {
    preload(src, {
      as: "image",
      fetchPriority: "high",
      imageSrcSet: srcSet,
      imageSizes: responsiveSizes,
    });
  }

  return (
    <styled.img
      alt={alt}
      loading={priority ? "eager" : "lazy"}
      decoding="async"
      fetchPriority={priority ? "high" : undefined}
      {...rest}
      src={src}
      srcSet={srcSet}
      sizes={responsiveSizes}
      htmlWidth={width}
      htmlHeight={height}
    />
  );
}
