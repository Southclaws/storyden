import { z } from "zod";

export const ThemeManifestSchema = z
  .object({
    css: z.array(z.string()).default([]),
    scripts: z.array(z.string()).default([]),
  })
  .default({
    css: [],
    scripts: [],
  });

export type ThemeManifest = z.infer<typeof ThemeManifestSchema>;

export const EMPTY_THEME_MANIFEST: ThemeManifest = {
  css: [],
  scripts: [],
};

export function parseThemeManifest(value: unknown): ThemeManifest {
  const parsed = ThemeManifestSchema.safeParse(value);
  if (!parsed.success) {
    return EMPTY_THEME_MANIFEST;
  }

  return {
    css: parsed.data.css,
    scripts: parsed.data.scripts,
  };
}

type FilterOptions = {
  webAddress: string;
  apiAddress: string;
};

export function filterAllowedThemeAssets(
  manifest: ThemeManifest,
  options: FilterOptions,
): ThemeManifest {
  const allowedOrigins = getAllowedOrigins(options);

  return {
    css: manifest.css
      .map((href) => resolveThemeAssetURL(href, options))
      .filter((href) => isAllowedThemeAssetURL(href, allowedOrigins)),
    scripts: manifest.scripts
      .map((src) => resolveThemeAssetURL(src, options))
      .filter((src) => isAllowedThemeAssetURL(src, allowedOrigins)),
  };
}

function isAllowedThemeAssetURL(value: string, allowedOrigins: Set<string>) {
  const trimmed = value.trim();
  if (!trimmed) {
    return false;
  }

  let parsed: URL;
  try {
    parsed = new URL(trimmed);
  } catch {
    return false;
  }

  return allowedOrigins.has(parsed.origin);
}

function resolveThemeAssetURL(value: string, options: FilterOptions) {
  const trimmed = value.trim();
  if (!trimmed) {
    return "";
  }

  if (!trimmed.startsWith("/")) {
    return trimmed;
  }

  const base = trimmed.startsWith("/api/")
    ? options.apiAddress
    : options.webAddress;

  try {
    return new URL(trimmed, base).toString();
  } catch {
    return "";
  }
}

function getAllowedOrigins(options: FilterOptions) {
  const values = [options.webAddress, options.apiAddress];

  const origins = new Set<string>();
  values.forEach((value) => {
    try {
      origins.add(new URL(value).origin);
    } catch {
      // Ignore malformed addresses and continue with valid origins.
    }
  });

  return origins;
}
