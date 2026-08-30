import "server-only";

import { unstable_cache } from "next/cache";

import { getAPIAddress, serverEnvironment } from "@/config";

import { EMPTY_THEME_MANIFEST, parseThemeManifest } from "./manifest";

export const THEME_CACHE_TAG = "theme";

const getCachedThemeManifest = unstable_cache(
  async () => fetchThemeManifest(),
  ["theme-manifest-v1"],
  { revalidate: 300, tags: [THEME_CACHE_TAG] },
);

export async function getServerThemeManifest() {
  return getCachedThemeManifest();
}

async function fetchThemeManifest() {
  try {
    const response = await fetch(`${getAPIAddress()}/api/info/theme`, {
      cache: "no-store",
      signal: AbortSignal.timeout(3_000),
    });
    if (!response.ok) {
      return EMPTY_THEME_MANIFEST;
    }

    return parseThemeManifest(
      await response.json(),
      serverEnvironment().API_ADDRESS,
    );
  } catch {
    return EMPTY_THEME_MANIFEST;
  }
}
