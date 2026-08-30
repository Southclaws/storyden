import { z } from "zod";

const ThemeAssetSchema = z.object({
  id: z.string(),
  filename: z.string(),
  path: z.string(),
  mime_type: z.enum(["text/css", "application/javascript"]),
  size: z.number().int().positive(),
  integrity: z.string().startsWith("sha256-"),
});

export const ThemeManifestSchema = z.object({
  api_version: z.literal("v1"),
  enabled: z.boolean(),
  revision: z.string(),
  stylesheets: z.array(ThemeAssetSchema).max(32),
  scripts: z.array(ThemeAssetSchema).max(32),
});

export type ThemeManifest = z.infer<typeof ThemeManifestSchema>;
export type ThemeAsset = z.infer<typeof ThemeAssetSchema>;
export type ResolvedThemeAsset = ThemeAsset & { href: string };
export type ResolvedThemeManifest = Omit<
  ThemeManifest,
  "stylesheets" | "scripts"
> & {
  stylesheets: ResolvedThemeAsset[];
  scripts: ResolvedThemeAsset[];
};

export const EMPTY_THEME_MANIFEST: ResolvedThemeManifest = {
  api_version: "v1",
  enabled: false,
  revision: "",
  stylesheets: [],
  scripts: [],
};

export function resolveThemeAssetHref(
  assetPath: string,
  publicAPIAddress: string,
): string | undefined {
  try {
    const apiURL = new URL(publicAPIAddress);
    const href = new URL(assetPath, apiURL);
    if (
      href.origin !== apiURL.origin ||
      !href.pathname.startsWith("/api/info/theme/assets/")
    ) {
      return undefined;
    }
    return href.toString();
  } catch {
    return undefined;
  }
}

export function parseThemeManifest(
  value: unknown,
  publicAPIAddress: string,
): ResolvedThemeManifest {
  const parsed = ThemeManifestSchema.safeParse(value);
  if (!parsed.success || !parsed.data.enabled) {
    return EMPTY_THEME_MANIFEST;
  }

  const resolve = (asset: ThemeAsset): ResolvedThemeAsset | undefined => {
    const href = resolveThemeAssetHref(asset.path, publicAPIAddress);
    return href ? { ...asset, href } : undefined;
  };

  const stylesheets = parsed.data.stylesheets
    .map(resolve)
    .filter((asset): asset is ResolvedThemeAsset => asset !== undefined);
  const scripts = parsed.data.scripts
    .map(resolve)
    .filter((asset): asset is ResolvedThemeAsset => asset !== undefined);

  if (
    stylesheets.length !== parsed.data.stylesheets.length ||
    scripts.length !== parsed.data.scripts.length
  ) {
    return EMPTY_THEME_MANIFEST;
  }

  return { ...parsed.data, stylesheets, scripts };
}
