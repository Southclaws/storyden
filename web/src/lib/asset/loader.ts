import type { ImageLoaderProps } from "next/image";

import { getImageVariantURL } from "@/utils/asset";

export default function storydenImageLoader({ src, width }: ImageLoaderProps) {
  const variant = width <= 320 ? "s" : width <= 960 ? "m" : "l";

  return getImageVariantURL(src, variant) ?? src;
}
