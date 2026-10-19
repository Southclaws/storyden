import type { Asset, AssetVariant } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";

export function getAssetURL(path?: string, variant?: AssetVariant) {
  if (!path) {
    return;
  }

  const source = /^(https?:|blob:|data:|\/\/)/.test(path)
    ? path
    : `${API_ADDRESS}${path}`;
  return variant ? getImageVariantURL(source, variant) : source;
}

export function getImageVariantURL(
  source: string | undefined,
  variant: AssetVariant,
) {
  if (!source) {
    return source;
  }

  try {
    const base = new URL(API_ADDRESS);
    const url = new URL(source, base);
    if (
      url.origin !== base.origin ||
      !url.pathname.startsWith("/api/assets/")
    ) {
      return source;
    }

    url.searchParams.set("variant", variant);
    return source.startsWith("/") && !source.startsWith("//")
      ? `${url.pathname}${url.search}${url.hash}`
      : url.toString();
  } catch {
    return source;
  }
}

export function getAssetImageProps(
  asset?: Pick<Asset, "path" | "width" | "height">,
  variant: Extract<AssetVariant, "s" | "m" | "l"> = "m",
) {
  const original = getAssetURL(asset?.path);
  const src = getImageVariantURL(original, variant);
  if (!src || src === original || !asset || !asset.width || !asset.height) {
    return { src };
  }

  const edge = Math.max(asset.width, asset.height);
  const sizes = { s: 320, m: 960, l: 1920 } as const;
  const scale = Math.min(1, sizes[variant] / edge);
  const seen = new Set<number>();
  const candidates = (["s", "m", "l"] as const).flatMap((preset) => {
    const width = Math.max(
      1,
      Math.floor(asset.width * Math.min(1, sizes[preset] / edge)),
    );
    if (seen.has(width)) {
      return [];
    }

    seen.add(width);

    return [`${getAssetURL(asset.path, preset)} ${width}w`];
  });

  return {
    src,
    srcSet: candidates.join(", "),
    width: Math.max(1, Math.floor(asset.width * scale)),
    height: Math.max(1, Math.floor(asset.height * scale)),
  };
}
